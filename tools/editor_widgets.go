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
	EDITOR_WIDTH  = 200
	EDITOR_HEIGHT = 30
)

type ModelWidgets struct {
	nodeWidgets       map[*model.Node]*NodeWidget
	connectionWidgets map[*model.Connection]*ConnectionWidget
	toolbar           *Toolbar
}

type NodeWidget struct {
	textBox    *widget.Editor
	boldButton *widget.Clickable
	textSizer  *TextSizeAdjuster
	isBold     bool
}

type ConnectionWidget struct {
	curveButton *widget.Clickable
}

type Toolbar struct {
	showGenerated *widget.Bool
	serifButton   *widget.Clickable
	sansButton    *widget.Clickable
	textSizer     *TextSizeAdjuster
}

type TextSizeAdjuster struct {
	increaseButton *widget.Clickable
	decreaseButton *widget.Clickable
}

func InitWidgets(m *model.Model) ModelWidgets {
	var w ModelWidgets

	w.nodeWidgets = make(map[*model.Node]*NodeWidget)
	w.connectionWidgets = make(map[*model.Connection]*ConnectionWidget)

	for _, c := range m.Connections {
		w.connectionWidgets[c] = &ConnectionWidget{
			curveButton: &widget.Clickable{},
		}
	}

	for _, n := range m.Nodes {
		nw := &NodeWidget{
			textBox:    new(widget.Editor),
			boldButton: new(widget.Clickable),
			textSizer: &TextSizeAdjuster{
				increaseButton: new(widget.Clickable),
				decreaseButton: new(widget.Clickable),
			},
		}
		nw.textBox.SingleLine = true
		w.nodeWidgets[n] = nw
	}

	w.toolbar = &Toolbar{
		showGenerated: new(widget.Bool),
		serifButton:   new(widget.Clickable),
		sansButton:    new(widget.Clickable),
		textSizer: &TextSizeAdjuster{
			increaseButton: new(widget.Clickable),
			decreaseButton: new(widget.Clickable),
		},
	}

	return w
}

func (w ModelWidgets) DrawNodeEditor(ops *op.Ops, gtx layout.Context, th *material.Theme, n *model.Node, pos utils.LocalPos, ec *EditContext, fontFaces []font.FontFace, isSerif bool, fontSize float32) {
	nodeWidget := w.nodeWidgets[n]

	globalPos := pos.ToGlobal(ec.scaleFactor, ec.viewportCenter, ec.windowSize)
	dim := utils.GlobalDim{W: EDITOR_WIDTH, H: EDITOR_HEIGHT}

	// Apply scale transform
	anchorPoint := globalPos.ToF32()
	defer op.Affine(f32.Affine2D{}.Scale(anchorPoint, f32.Pt(ec.scaleFactor, ec.scaleFactor))).Push(ops).Pop()
	defer op.Offset(globalPos.SubDim(dim.Div(2)).ToImagePnt()).Push(ops).Pop()

	// Draw background
	radius := 2
	utils.DrawRoundedRect(ops,
		utils.GlobalPos{X: EDITOR_WIDTH / 2, Y: EDITOR_HEIGHT / 2},
		dim,
		radius,
		color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		2)

	editBarGtx := gtx
	editBarGtx.Constraints = layout.Exact(image.Pt(dim.W, dim.H))

	// Layout contents
	layout.UniformInset(unit.Dp(2.5)).Layout(editBarGtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Square bold button on the left
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if nodeWidget.boldButton.Clicked(gtx) {
					nodeWidget.isBold = !nodeWidget.isBold
				}

				btn := material.Button(th, nodeWidget.boldButton, "B")
				btn.Background = color.NRGBA{R: 100, G: 100, B: 200, A: 255}
				btn.TextSize = unit.Sp(12)
				if nodeWidget.isBold {
					btn.Font.Weight = font.Weight(300)
				}
				btn.Inset = layout.Inset{Top: 1, Bottom: 1, Left: 1, Right: 1}
				btn.CornerRadius = unit.Dp(2)
				if nodeWidget.isBold {
					btn.Background = color.NRGBA{R: 100, G: 100, B: 170, A: 255}
				}

				// Make button square by constraining to fixed size
				buttonSize := float32(EDITOR_HEIGHT - 5)
				gtx.Constraints = layout.Exact(image.Pt(int(buttonSize), int(buttonSize)))

				return btn.Layout(gtx)
			}),

			// Spacer
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

			// Text editor takes remaining space
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, nodeWidget.textBox, n.Text)
				editor.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
				if nodeWidget.isBold {
					editor.Font.Weight = font.Bold
				} else {
					editor.Font.Weight = font.Normal
				}
				editor.TextSize = unit.Sp(14)
				return editor.Layout(gtx)
			}),
		)
	})

	fontFace := utils.GetFontFace(isSerif, nodeWidget.isBold, fontFaces)

	if n.Bold != nodeWidget.isBold {
		n.Bold = nodeWidget.isBold
		n.TextWidth = utils.GetTextWidth(n.Text, fontFace, fontSize, gtx)
	}

	newText := nodeWidget.textBox.Text()
	if n.Text != newText && newText != "" {
		n.Text = newText
		n.TextWidth = utils.GetTextWidth(n.Text, fontFace, fontSize, gtx)
	}

	if newText == "" && len(n.Text) == 1 {
		n.Text = n.VarName
		n.TextWidth = utils.GetTextWidth(n.Text, fontFace, fontSize, gtx)
	} else if n.Text != n.VarName && newText == "" { // case when loading a model that already has edited text
		nodeWidget.textBox.SetText(n.Text)
		newText = n.Text
	}
}

func (w ModelWidgets) DrawToolbar(ops *op.Ops, gtx layout.Context, th *material.Theme, m *model.Model) {
	tb := w.toolbar
	height := 50

	tbGtx := gtx
	tbGtx.Constraints.Max = image.Point{X: gtx.Constraints.Max.X, Y: height}

	// draw background
	dim := utils.GlobalDim{
		W: gtx.Constraints.Min.X,
		H: height,
	}
	utils.DrawRect(ops, utils.GlobalPos{}.AddDim(dim.Div(2)), dim, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, 1)

	// Layout contents
	layout.UniformInset(unit.Dp(2.5)).Layout(tbGtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// draw checkbox for showing generated connections
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					checkBox := material.CheckBox(th, tb.showGenerated, "Show generated connections")
					checkBox.TextSize = unit.Sp(14)
					return checkBox.Layout(gtx)
				}),

				// spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),

				// draw font family selection
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					serifBtn := material.Button(th, tb.serifButton, "Serif")
					serifBtn.Font.Typeface = m.Font.Faces[2].Font.Typeface
					serifBtn.TextSize = unit.Sp(16)
					serifBtn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
					serifBtn.CornerRadius = unit.Dp(2)
					if m.Font.IsSerif {
						serifBtn.Background = color.NRGBA{A: 100}
					} else {
						serifBtn.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					}
					return serifBtn.Layout(gtx)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					sansBtn := material.Button(th, tb.sansButton, "Sans Serif")
					sansBtn.Font.Typeface = m.Font.Faces[0].Font.Typeface
					sansBtn.TextSize = unit.Sp(16)
					sansBtn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
					sansBtn.CornerRadius = unit.Dp(2)
					if !m.Font.IsSerif {
						sansBtn.Background = color.NRGBA{A: 100}
					} else {
						sansBtn.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					}
					return sansBtn.Layout(gtx)
				}),

				// spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),

				// draw font size selector
			)
		})
	})

	if tb.sansButton.Pressed() {
		m.Font.IsSerif = false
		model.RecalcTextWidths(m, gtx)
	}

	if tb.serifButton.Pressed() {
		m.Font.IsSerif = true
		model.RecalcTextWidths(m, gtx)
	}

	if tb.showGenerated.Pressed() {
		m.ViewGenerated = !m.ViewGenerated
	}
}
