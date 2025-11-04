package utils

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func MakeRect(pos GlobalPos, dim GlobalDim) image.Rectangle {
	minPt := image.Point{
		X: pos.X - dim.W/2,
		Y: pos.Y - dim.H/2,
	}
	maxPt := image.Point{
		X: pos.X + dim.W/2,
		Y: pos.Y + dim.H/2,
	}

	return image.Rectangle{Min: minPt, Max: maxPt}
}

func DrawRect(ops *op.Ops, pos GlobalPos, dim GlobalDim, col color.NRGBA, thickness float32) {
	rect := clip.Rect(MakeRect(pos, dim))

	defer rect.Push(ops).Pop()

	// Draw fill
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)

	// Draw outline
	paint.FillShape(ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		clip.Stroke{
			Path:  rect.Path(),
			Width: thickness,
		}.Op(),
	)
}

func DrawEllipse(ops *op.Ops, pos GlobalPos, dim GlobalDim, col color.NRGBA, thickness float32) {
	el := clip.Ellipse(MakeRect(pos, dim))

	defer el.Push(ops).Pop()

	// Draw fill
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)

	// Draw outline
	paint.FillShape(ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		clip.Stroke{
			Path:  el.Path(ops),
			Width: thickness,
		}.Op(),
	)
}

func DrawLine(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness float32) {
	var path clip.Path
	path.Begin(ops)
	path.MoveTo(posA.ToF32())
	path.LineTo(posB.ToF32())
	path.Close()

	paint.FillShape(ops, col,
		clip.Stroke{
			Path:  path.End(),
			Width: thickness,
		}.Op(),
	)
}
