package main

import (
	"main/model"
	"main/utils"
	"math"
	"slices"
	"sort"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
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

func DrawModel(ops *op.Ops, gtx layout.Context, m *model.Model, ec *EditContext, selectedNode *model.Node) {
	// Reset all node connections every frame
	for _, n := range m.Nodes {
		if !n.Visible {
			continue
		}

		n.EdgeConnections = [4][]*model.Connection{}
		// define node dimensions
		if n.TextWidth == 0 {
			n.TextWidth = utils.GetTextWidth(n.Text, m.Font.Face, m.Font.Size)
		}
		// todo: decide whether to snap dimensions to grid as well as position
		//adjWidth := utils.SnapValue(textWidth+targetPadding*2, ec.snapGridSize)
		adjWidth := n.TextWidth + targetPadding*2
		n.Padding = (adjWidth - n.TextWidth) / 2.0
		switch n.Class {
		case model.OBSERVED:
			n.Dim = utils.LocalDim{W: adjWidth, H: 50}
		case model.LATENT:
			n.Dim = utils.LocalDim{W: adjWidth, H: adjWidth}
		case model.INTERCEPT:
			//todo: handle intercepts
		}
	}

	// Rule for connections:

	// for observed variables --
	// regression arrows originate and terminate at the center of an edge, unless there are more than one regression arrows
	// from the same edge, in which case the terminating location is evenly distributed. The only exception is if
	// the connection is at a cardinal angle, in which case it originates/terminates from the center, regardless
	// of how many other connections originate/terminate at the same node.

	// for latent variables --
	// all arrows terminating at a latent variable terminate at the place along the edge that makes the path
	// the shortest.
	// all arrows originating from a latent variable originate at the place along the edge that makes the path
	// the shortest

	for _, c := range m.Connections {
		if !c.UserDefined && !m.ViewGenerated {
			continue
		}

		switch c.Type {
		case model.STRAIGHT:
			c.Angle = utils.GetAngleLoc(c.Origin.Pos, c.Destination.Pos)
		case model.CURVED:
			ctrl := utils.GetCtrlPoint(c.Origin.Pos.ToF32(), c.Destination.Pos.ToF32(), c.Curvature)
			angle := -math.Atan2(float64(ctrl.Y-c.Origin.Pos.ToF32().Y), float64(ctrl.X-c.Origin.Pos.ToF32().X))
			c.Angle = utils.NormalizeAngle(angle)
		}

		AssignToEdges(c, m.Nodes)
	}

	// calculate observed nodes
	// must handle observed before latent to establish positions of connection ends
	for _, n := range m.Nodes {
		if !n.Visible {
			continue
		}

		switch n.Class {
		case model.OBSERVED:
			// handle connections
			for e := 0; e < 4; e++ {
				connections := n.EdgeConnections[e]

				if len(connections) == 0 {
					continue
				}

				// sort connections based on angle
				angleMap := make(map[*model.Connection]float64, len(connections))
				for _, c := range connections {
					switch c.Type {
					case model.STRAIGHT:
						if c.Destination == n {
							angleMap[c] = c.Angle
						}
						if c.Origin == n {
							angleMap[c] = utils.NormalizeAngle(c.Angle + math.Pi)
						}
					case model.CURVED:
						if c.Origin == n {
							angleMap[c] = utils.GetAngleLoc(c.Destination.Pos, c.Origin.Pos)
						}
						if c.Destination == n {
							angleMap[c] = utils.GetAngleLoc(c.Origin.Pos, c.Destination.Pos)
						}
					}
				}

				// do the sorting
				// edge 3 is a special case, since the angle values wrap around from 2Pi to 0
				if e == 3 {
					for k, v := range angleMap {
						if v < math.Pi {
							angleMap[k] += 2 * math.Pi
						}
					}
				}

				sort.Slice(connections, func(i, j int) bool {
					return angleMap[connections[i]] < angleMap[connections[j]]
				})

				edgePoints := SubdivideNodeEdge(n, e, len(connections))

				for i, c := range connections {
					// override edge offset positioning if the nodes are sufficiently in-line
					// e.g., if two nodes are directly on top of each other, they should connect from the middle,
					// no matter the number of other connections
					// this only matters if the number of connections on an edge is even
					// for now, I will only consider the problem when the number of connections is 2
					angle := utils.GetAngleLoc(c.Origin.Pos, c.Destination.Pos)
					if len(connections) == 2 && utils.SufficientlyAligned(angle, math.Pi/256) {
						if c.Destination == n {
							c.DestinationPos = SubdivideNodeEdge(c.Destination, e, 1)[0]
						}
						if c.Origin == n {
							c.OriginPos = SubdivideNodeEdge(c.Origin, e, 1)[0]
						}
						continue
					}

					if c.Destination == n {
						c.DestinationPos = edgePoints[i]
					}
					if c.Origin == n {
						c.OriginPos = edgePoints[i]
					}
				}
			}
		}
	}

	for _, n := range m.Nodes {
		if !n.Visible {
			continue
		}

		switch n.Class {
		case model.OBSERVED:
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

		textOffset := utils.LocalDim{W: n.Dim.W/2.0 - n.Padding, H: m.Font.Size / 1.5} // I think 1.5 is a magic number
		utils.DrawText(
			ops,
			gtx,
			n.Pos.SubDim(textOffset).ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
			n.Text,
			m.Font.Face,
			unit.Sp(m.Font.Size),
			ec.scaleFactor,
		)
	}

	for _, c := range m.Connections {
		if !c.UserDefined && !m.ViewGenerated {
			continue
		}

		if c.EstWidth == 0 {
			c.EstText, c.EstDim, c.EstWidth = utils.CalculateEstimate(m.Font.Face, m.Font.Size-2, m.CoeffDisplay, c.Est, c.PValue, c.CI, 2, c.EstPadding)
		}

		if c.Origin.Class == model.LATENT {
			angleFromLatent := utils.GetAngleLoc(c.Origin.Pos, c.DestinationPos)
			c.OriginPos = utils.MoveAlongAngleLoc(c.Origin.Pos, angleFromLatent, c.Origin.Dim.W/2.0)
		}
		if c.Destination.Class == model.LATENT {
			angleToLatent := utils.GetAngleLoc(c.OriginPos, c.Destination.Pos)
			c.DestinationPos = utils.MoveAlongAngleLoc(c.Destination.Pos, angleToLatent+math.Pi, c.Destination.Dim.W/2.0)
		}

		switch c.Type {
		case model.STRAIGHT:
			utils.DrawArrowLine(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				ec.windowSize,
			)

			// determine label position as distance along curve
			angle := utils.GetAngleLoc(c.OriginPos, c.DestinationPos)
			dist := utils.DistLoc(c.OriginPos, c.DestinationPos)
			c.EstPos = utils.MoveAlongAngleLoc(c.OriginPos, angle, dist*c.AlongLineProp)
		case model.CURVED:
			utils.DrawArrowArc(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				c.Curvature,
				ec.windowSize,
			)

			// determine label position as distance along curve
			ctrl := utils.GetCtrlPoint(c.OriginPos.ToF32(), c.DestinationPos.ToF32(), c.Curvature)
			c.EstPos = utils.MoveAlongBezier(c.OriginPos.ToF32(), c.DestinationPos.ToF32(), ctrl, c.AlongLineProp)
		}
	}

	// draw estimate labels after ALL of the connections to ensure proper layering
	for _, c := range m.Connections {
		if !c.UserDefined && !m.ViewGenerated {
			continue
		}

		if m.CoeffDisplay != utils.NONE {
			utils.DrawEstimate(
				ops,
				gtx,
				c.EstPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				m.Font.Face,
				m.Font.Size,
				ec.scaleFactor,
				c.EstPadding,
				c.EstText,
				c.EstDim,
				c.EstWidth,
			)
		}
	}
}

// DrawModelFixed is a faster function for lazyUpdate the view. It does not update local positions
func DrawModelFixed(ops *op.Ops, gtx layout.Context, m *model.Model, ec *EditContext) {
	for _, n := range m.Nodes {
		if !n.Visible {
			continue
		}

		switch n.Class {
		case model.OBSERVED:
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
		}

		textOffset := utils.LocalDim{W: n.Dim.W/2.0 - n.Padding, H: m.Font.Size / 1.5} // I think 1.5 is a magic number
		utils.DrawText(
			ops,
			gtx,
			n.Pos.SubDim(textOffset).ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
			n.Text,
			m.Font.Face,
			unit.Sp(m.Font.Size),
			ec.scaleFactor,
		)
	}

	for _, c := range m.Connections {
		if !c.UserDefined && !m.ViewGenerated {
			continue
		}

		switch c.Type {
		case model.STRAIGHT:
			utils.DrawArrowLine(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				ec.windowSize,
			)

			// determine label position as distance along curve
			angle := utils.GetAngleLoc(c.OriginPos, c.DestinationPos)
			dist := utils.DistLoc(c.OriginPos, c.DestinationPos)
			c.EstPos = utils.MoveAlongAngleLoc(c.OriginPos, angle, dist*c.AlongLineProp)
		case model.CURVED:
			utils.DrawArrowArc(
				ops,
				c.OriginPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.DestinationPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				c.Col,
				c.Thickness*ec.scaleFactor,
				c.Curvature,
				ec.windowSize,
			)

			// determine label position as distance along curve
			ctrl := utils.GetCtrlPoint(c.OriginPos.ToF32(), c.DestinationPos.ToF32(), c.Curvature)
			c.EstPos = utils.MoveAlongBezier(c.OriginPos.ToF32(), c.DestinationPos.ToF32(), ctrl, c.AlongLineProp)
		}
	}
	// draw estimate labels after ALL of the connections to ensure proper layering
	for _, c := range m.Connections {
		if !c.UserDefined && !m.ViewGenerated {
			continue
		}

		if m.CoeffDisplay != utils.NONE {
			utils.DrawEstimate(
				ops,
				gtx,
				c.EstPos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
				m.Font.Face,
				m.Font.Size,
				ec.scaleFactor,
				c.EstPadding,
				c.EstText,
				c.EstDim,
				c.EstWidth,
			)
		}
	}
}

func AssignToEdges(c *model.Connection, nodes []*model.Node) {
	switch {
	case c.Origin.Class == model.OBSERVED && c.Destination.Class == model.LATENT:
		candidateOriginEdges := GetCandidateDestEdges(InvertAngle(c))
		edgeOrigin := GetBestEdge(candidateOriginEdges, c.Origin, c.Destination, c.Destination.Pos, nodes)
		c.Origin.EdgeConnections[edgeOrigin] = append(c.Origin.EdgeConnections[edgeOrigin], c)

	case c.Origin.Class == model.LATENT && c.Destination.Class == model.OBSERVED:
		candidateOriginEdges := GetCandidateDestEdges(c.Angle)
		edgeDest := GetBestEdge(candidateOriginEdges, c.Destination, c.Origin, c.Origin.Pos, nodes)
		c.Destination.EdgeConnections[edgeDest] = append(c.Destination.EdgeConnections[edgeDest], c)

	case c.Origin.Class == model.OBSERVED && c.Destination.Class == model.OBSERVED:
		// assign origin edge like normal
		edgeOrigin := AngleToEdge(c.Angle)
		c.Origin.EdgeConnections[edgeOrigin] = append(c.Origin.EdgeConnections[edgeOrigin], c)

		// from this edge, find recalculate the connection angle
		edgePoint := SubdivideNodeEdge(c.Origin, edgeOrigin, 1)[0]
		angleIntermediate := utils.GetAngleLoc(edgePoint, c.Destination.Pos)

		candidateDestEdges := GetCandidateDestEdges(angleIntermediate)
		edgeDest := GetBestEdge(candidateDestEdges, c.Destination, c.Origin, edgePoint, nodes)
		c.Destination.EdgeConnections[edgeDest] = append(c.Destination.EdgeConnections[edgeDest], c)
	}

}

func GetBestEdge(candidateEdges []int, destNode, originNode *model.Node, originEdgePoint utils.LocalPos, nodes []*model.Node) int {
	// best edge is defined as the edge where the line angle is maximally different from
	// the edge angle (the edge angle being either flat or 90* based on the edge)
	// OR
	// the edge that does not intersect another node

	intersections := make([]bool, 2)
	angles := make([]float64, 2)
	for i, edge := range candidateEdges {
		destEdgePoint := SubdivideNodeEdge(destNode, edge, 1)[0]
		angle := utils.GetAngleLoc(originEdgePoint, destEdgePoint)
		switch {
		case edge == 0 || edge == 2:
			angles[i] = math.Abs(math.Sin(angle))
		case edge == 1 || edge == 3:
			angles[i] = math.Abs(math.Cos(angle))
		}

		intersections[i] = checkIntersection(originEdgePoint, destEdgePoint, nodes, []*model.Node{destNode, originNode})
	}

	switch {
	case intersections[0] && !intersections[1]:
		return candidateEdges[1]
	case !intersections[0] && intersections[1]:
		return candidateEdges[0]
	}

	// find the largest difference
	sort.Slice(candidateEdges, func(i, j int) bool {
		return angles[i] > angles[j]
	})

	return candidateEdges[0]
}

func GetCandidateDestEdges(angle float64) []int {
	switch {
	case angle >= 0 && angle < math.Pi/2.0:
		return []int{2, 3}
	case angle >= math.Pi/2.0 && angle < math.Pi:
		return []int{1, 2}
	case angle >= math.Pi && angle < 3*math.Pi/2.0:
		return []int{0, 1}
	case angle >= 3*math.Pi/2.0 && angle < 2*math.Pi:
		return []int{3, 0}
	default:
		panic("invalid angle value")
	}
}

func InvertAngle(c *model.Connection) float64 {
	var angle float64
	switch c.Type {
	case model.STRAIGHT:
		angle = utils.NormalizeAngle(c.Angle + math.Pi)
	case model.CURVED:
		ctrl := utils.GetCtrlPoint(c.Origin.Pos.ToF32(), c.Destination.Pos.ToF32(), c.Curvature)
		angle = -math.Atan2(float64(ctrl.Y-c.Destination.Pos.ToF32().Y), float64(ctrl.X-c.Destination.Pos.ToF32().X))
		angle = utils.NormalizeAngle(angle)
	}
	return angle
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

func checkIntersection(originPos, destPos utils.LocalPos, nodes, exclusion []*model.Node) bool {
	for _, n := range nodes {
		if slices.Contains(exclusion, n) {
			continue
		}

		nodeRect := utils.LocalRect{
			NW: n.Pos.SubDim(n.Dim.Div(2)),
			SE: n.Pos.AddDim(n.Dim.Div(2)),
		}

		if utils.SegmentIntersectsRect(originPos, destPos, nodeRect) {
			return true
		}
	}
	return false
}
