package main

import (
	"main/model"
	"main/utils"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
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
		w.nodeWidgets[n] = new(NodeWidget)
	}

	return w
}

func (w ModelWidgets) DrawNodeEditor(ops *op.Ops, gtx layout.Context, th *material.Theme, n *model.Node, pos utils.LocalPos, ec *EditContext) {
	//nodeWidget := w.nodeWidgets[n]
	//
	//// Convert position to global coordinates
	//globalPos := pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize)
	//
	//// Create a constrained context for the toolbar dimensions
	//dim := utils.LocalDim{
	//	W: EDITOR_WIDTH,
	//	H: EDITOR_HEIGHT,
	//}
	//globalDim := dim.ToGlobal(ec.scaleFactor)
	//
	//// Stack to position the toolbar at the specified location
	//stack := op.Offset(globalPos.ToImagePnt()).Push(ops)
	//
	//// Create a constrained layout context
	//toolbarGtx := gtx
	//toolbarGtx.Constraints = layout.Exact(image.Pt(int(globalDim.W), int(globalDim.H)))
	//
	//// Draw background with rounded corners
	//radius := float32(10 * ec.scaleFactor)
	//
	//// Background macro to draw behind everything
	//macro := op.Record(ops)
	//utils.DrawRoundedRect(ops,
	//	utils.GlobalPos{X: 0, Y: 0}, // Draw at origin since we're already offset
	//	globalDim,
	//	int(radius),
	//	color.NRGBA{R: 40, G: 40, B: 50, A: 100}, // Solid background
	//	0)
	//callOps := macro.Stop()
	//callOps.Add(ops)
	//
	//// Layout the toolbar contents
	//layout.Flex{
	//	Axis: layout.Vertical,
	//}.Layout(toolbarGtx,
	//	// Text editor takes most of the space
	//	layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
	//		// Add padding around the editor
	//		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//			// Create a styled editor
	//			editor := material.Editor(th, &nodeWidget.textBox, "Type here...")
	//			editor.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	//
	//			// Apply bold font if needed
	//			if nodeWidget.isBold {
	//				editor.Font.Weight = font.Bold
	//			} else {
	//				editor.Font.Weight = font.Normal
	//			}
	//
	//			return editor.Layout(gtx)
	//		})
	//	}),
	//
	//	// Button row at the bottom
	//	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//			// Check if button was clicked
	//			if nodeWidget.boldButton.Clicked(gtx) {
	//				nodeWidget.isBold = !nodeWidget.isBold
	//			}
	//
	//			// Create button with dynamic text
	//			buttonText := "Bold"
	//			if nodeWidget.isBold {
	//				buttonText = "Unbold"
	//			}
	//
	//			btn := material.Button(th, &nodeWidget.boldButton, buttonText)
	//			btn.Background = color.NRGBA{R: 70, G: 130, B: 180, A: 255}
	//
	//			return btn.Layout(gtx)
	//		})
	//	}),
	//)
	//
	//stack.Pop()
}
