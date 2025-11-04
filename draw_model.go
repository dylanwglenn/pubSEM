package main

import (
	"main/model"
	"main/utils"

	"gioui.org/op"
)

func DrawModel(ops *op.Ops, m *model.Model, ec *EditContext) {

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
		utils.DrawLine(
			ops,
			c.Origin.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
			c.Destination.Pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
			c.Col,
			c.Thickness*ec.scaleFactor,
		)
	}

	// draw nodes
	for _, n := range m.Nodes {
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
	}
}
