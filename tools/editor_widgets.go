package main

import (
	"image"
	"image/color"
	"main/model"
	"main/utils"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const (
	EDITOR_WIDTH  = 175
	EDITOR_HEIGHT = 25
)

type ModelWidgets struct {
	nodeWidgets       map[*model.Node]*NodeWidget
	connectionWidgets map[*model.Connection]*ConnectionWidget
}

type NodeWidget struct {
	textBox    widget.Editor
	boldButton widget.Clickable
	isBold     bool
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
		nw := new(NodeWidget)
		nw.textBox.SingleLine = true
		w.nodeWidgets[n] = nw
	}

	return w
}

func (w ModelWidgets) DrawNodeEditor(ops *op.Ops, gtx layout.Context, th *material.Theme, n *model.Node, pos utils.LocalPos, ec *EditContext) {
	nodeWidget := w.nodeWidgets[n]

	globalPos := pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize)
	dim := utils.LocalDim{W: EDITOR_WIDTH, H: EDITOR_HEIGHT}
	globalDim := dim.ToGlobal(ec.scaleFactor)

	// Apply scale transform
	anchorPoint := f32.Point{
		float32(globalPos.X),
		float32(globalPos.Y + globalDim.H/2),
	}
	defer op.Affine(f32.Affine2D{}.Scale(anchorPoint, f32.Pt(ec.scaleFactor, ec.scaleFactor))).Push(ops).Pop()
	// Stack to position the toolbar at the specified location
	defer op.Offset(globalPos.SubDim(globalDim.Div(2)).ToImagePnt()).Push(ops).Pop()

	toolbarGtx := gtx
	toolbarGtx.Constraints = layout.Exact(image.Pt(globalDim.W, globalDim.H))

	// Draw background
	radius := 2 * ec.scaleFactor
	utils.DrawRoundedRect(ops,
		utils.GlobalPos{X: globalDim.W / 2, Y: globalDim.H / 2},
		globalDim,
		int(radius),
		color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		2)

	// Layout contents
	layout.UniformInset(unit.Dp(10/ec.scaleFactor)).Layout(toolbarGtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Square bold button on the left
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if nodeWidget.boldButton.Clicked(gtx) {
					nodeWidget.isBold = !nodeWidget.isBold
				}

				btn := material.Button(th, &nodeWidget.boldButton, "B")
				btn.Background = color.NRGBA{R: 70, G: 130, B: 180, A: 255}
				btn.TextSize = unit.Sp(10)
				if nodeWidget.isBold {
					btn.Background = color.NRGBA{R: 100, G: 160, B: 210, A: 255}
				}

				// Make button square by constraining to fixed size
				buttonSize := float32(EDITOR_HEIGHT) * ec.scaleFactor
				gtx.Constraints = layout.Exact(image.Pt(int(buttonSize), int(buttonSize)))

				return btn.Layout(gtx)
			}),

			// Spacer
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

			// Text editor takes remaining space
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, &nodeWidget.textBox, n.Text)
				editor.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
				if nodeWidget.isBold {
					editor.Font.Weight = font.Bold
				} else {
					editor.Font.Weight = font.Normal
				}
				editor.TextSize = unit.Sp(10 * ec.scaleFactor)
				return editor.Layout(gtx)
			}),
		)
	})
}
