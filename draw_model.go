package main

import (
	"main/model"
	"main/utils"
	"math"
	"sort"

	"gioui.org/op"
)

// edge numbering for rectangles
//                0
//   ----------------------------
//   |                          |
//   |                          |
// 3 |                          | 1
//   |                          |
//   |                          |
//   ----------------------------
//                2

func DrawModel(ops *op.Ops, m *model.Model, ec *EditContext) {

	// Reset all node connections every frame
	for _, n := range m.Nodes {
		n.EdgeConnections = [4][]*model.Connection{}
	}

	// Rule for connections:

	// for observed variables --
	// regression arrows terminate at the center of an edge, unless there are more than one regression arrows
	// terminating on the same edge, in which case the terminating location is evenly distributed.
	// regression arrows should originate from the center of an edge, regardless of how many other arrows
	// originate from the same edge

	// for latent variables --
	// all arrows terminating at a latent variable terminate at the place along the edge that makes the path
	// the shortest.
	// all arrows originating from a latent variable originate at the place along the edge that makes the path
	// the shortest

	for _, c := range m.Connections {
		switch c.Type {
		case model.REGRESSION:
			c.Angle = utils.GetAngleLoc(c.Origin.Pos, c.Destination.Pos)
		case model.COVARIANCE:
			ctrl := utils.GetCtrlPoint(c.Origin.Pos.ToF32(), c.Destination.Pos.ToF32(), roundness, c.Curvature)
			angle := -math.Atan2(float64(ctrl.Y-c.Origin.Pos.ToF32().Y), float64(ctrl.X-c.Origin.Pos.ToF32().X))
			c.Angle = utils.NormalizeAngle(angle)
		}

		if c.Origin.Class == model.OBSERVED {
			edge := AngleToEdge(c.Angle)
			c.Origin.EdgeConnections[edge] = append(c.Origin.EdgeConnections[edge], c)
		}

		if c.Destination.Class == model.OBSERVED {
			var angle float64
			switch c.Type {
			case model.REGRESSION:
				angle = utils.NormalizeAngle(c.Angle + math.Pi)
			case model.COVARIANCE:
				ctrl := utils.GetCtrlPoint(c.Origin.Pos.ToF32(), c.Destination.Pos.ToF32(), roundness, c.Curvature)
				angle = -math.Atan2(float64(ctrl.Y-c.Destination.Pos.ToF32().Y), float64(ctrl.X-c.Destination.Pos.ToF32().X))
				angle = utils.NormalizeAngle(angle)
			}

			edge := AngleToEdge(angle)
			c.Destination.EdgeConnections[edge] = append(c.Destination.EdgeConnections[edge], c)
		}
	}

	// draw observed nodes
	// must handle observed before latent to establish positions of connection ends
	for _, n := range m.Nodes {
		switch n.Class {
		case model.OBSERVED:
			// handle connections
			for e := 0; e < 4; e++ {
				connections := n.EdgeConnections[e]

				if len(connections) == 0 {
					continue
				}

				// sort connections based on angle
				angles := make([]float64, len(connections))
				for i, c := range connections {
					if c.Destination == n {
						angles[i] = c.Angle
					}
					if c.Origin == n {
						angles[i] = utils.NormalizeAngle(c.Angle + math.Pi)
					}
				}

				// do the sorting
				// edge 3 is a special case, since the angle values wrap around from 2Pi to 0
				switch {
				case e == 3:
					for i := range angles {
						if angles[i] < math.Pi {
							angles[i] += 2 * math.Pi
						}
					}
					sort.Slice(connections, func(i, j int) bool {
						return angles[i] < angles[j]
					})
				default:
					sort.Slice(connections, func(i, j int) bool {
						return angles[i] < angles[j]
					})
				}

				edgePoints := SubdivideNodeEdge(n, e, len(connections))

				for i, c := range connections {
					if c.Destination == n {
						c.DestinationPos = edgePoints[i]
					}
					if c.Origin == n {
						c.OriginPos = edgePoints[i]
					}
				}
			}

			// draw the node itself
			utils.DrawRect(
				ops,
				n.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				n.Dim.ToGlobal(ec.scaleFactor),
				n.Col,
				n.Thickness*ec.scaleFactor,
			)

		case model.LATENT:
			utils.DrawEllipse(
				ops,
				n.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				n.Dim.ToGlobal(ec.scaleFactor),
				n.Col,
				n.Thickness*ec.scaleFactor,
			)

		case model.INTERCEPT:
			//TODO: Handle drawing intercept
		}
	}

	for _, c := range m.Connections {
		if c.Origin.Class == model.LATENT {
			c.OriginPos = utils.MoveAlongAngleLoc(c.Origin.Pos, c.Angle, c.Origin.Dim.W/2.0)
		}
		if c.Destination.Class == model.LATENT {
			c.DestinationPos = utils.MoveAlongAngleLoc(c.Destination.Pos, c.Angle+math.Pi, c.Destination.Dim.W/2.0)
		}

		switch c.Type {
		case model.REGRESSION:
			utils.DrawArrowLine(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				ec.windowSize,
			)
		case model.COVARIANCE:
			utils.DrawArrowArc(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				roundness,
				c.Curvature,
				ec.windowSize,
			)
		}
	}
}

func AngleToEdge(angle float64) int {
	switch {
	case angle >= math.Pi/4.0 && angle < 3*math.Pi/4.0:
		return 0
	case angle >= 3*math.Pi/4.0 && angle < 5*math.Pi/4.0:
		return 3
	case angle >= 5*math.Pi/4.0 && angle < 7*math.Pi/4.0:
		return 2
	case angle >= 7*math.Pi/4.0 || angle < math.Pi/4.0:
		return 1
	default:
		return -1
	}
}

func SubdivideNodeEdge(n *model.Node, edge, numPoints int) []utils.LocalPos {
	if n.Class != model.OBSERVED {
		return nil
	}

	if numPoints <= 0 {
		return nil
	}

	// initialize result slice
	res := make([]utils.LocalPos, numPoints)

	switch {
	case edge == 0:
		step := n.Dim.W / float32(numPoints+1)
		for i := 0; i < numPoints; i++ {
			res[i] = utils.LocalPos{
				X: (n.Pos.X + n.Dim.W/2.0) - step*float32(i+1),
				Y: n.Pos.Y - n.Dim.H/2.0,
			}
		}
	case edge == 2:
		step := n.Dim.W / float32(numPoints+1)
		for i := 0; i < numPoints; i++ {
			res[i] = utils.LocalPos{
				X: (n.Pos.X - n.Dim.W/2.0) + step*float32(i+1),
				Y: n.Pos.Y + n.Dim.H/2.0,
			}
		}
	case edge == 1:
		step := n.Dim.H / float32(numPoints+1)
		for i := 0; i < numPoints; i++ {
			res[i] = utils.LocalPos{
				X: n.Pos.X + n.Dim.W/2.0,
				Y: (n.Pos.Y + n.Dim.H/2.0) - step*float32(i+1),
			}
		}
	case edge == 3:
		step := n.Dim.H / float32(numPoints+1)
		for i := 0; i < numPoints; i++ {
			res[i] = utils.LocalPos{
				X: n.Pos.X - n.Dim.W/2.0,
				Y: (n.Pos.Y - n.Dim.H/2.0) + step*float32(i+1),
			}
		}
	default:
		panic("invalid edge")
	}

	return res
}
