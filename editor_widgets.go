package main

import (
	"image/color"
	"main/model"
	"main/utils"

	"gioui.org/op"
	"gioui.org/widget"
)

const (
	EDITOR_WIDTH  = 250
	EDITOR_HEIGHT = 40
)

type ModelWidgets struct {
	nodeWidgets       map[*model.Node]*NodeWidget
	connectionWidgets map[*model.Connection]*ConnectionWidget
}

type NodeWidget struct {
	textBox    widget.Editor
	boldButton widget.Clickable
}

type ConnectionWidget struct {
	curveButton widget.Clickable
}

func InitWidgets(m *model.Model) ModelWidgets {
	var w ModelWidgets

	w.nodeWidgets = make(map[*model.Node]*NodeWidget)
	w.connectionWidgets = make(map[*model.Connection]*ConnectionWidget)

	for _, c := range m.Connections {
		w.connectionWidgets[c] = new(ConnectionWidget)
	}

	for _, n := range m.Nodes {
		w.nodeWidgets[n] = new(NodeWidget)
	}

	return w
}

func (w ModelWidgets) DrawNodeEditor(ops *op.Ops, n *model.Node, pos utils.LocalPos, ec *EditContext) {
	// draw rounded rectangle background
	dim := utils.LocalDim{
		W: EDITOR_WIDTH,
		H: EDITOR_HEIGHT,
	}
	radius := 10 * ec.scaleFactor
	utils.DrawRoundedRect(ops,
		pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize),
		dim.ToGlobal(ec.scaleFactor),
		int(radius),
		color.NRGBA{A: 50},
		0)

}
