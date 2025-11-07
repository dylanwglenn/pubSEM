package utils

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
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

func DrawArrowArc(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness, roundness float32, curvature bool, windowSize GlobalDim) {
	// Calculate control point for tangent angles
	ctrl := GetCtrlPoint(posA.ToF32(), posB.ToF32(), roundness, curvature)

	arrowSize := float64(thickness * 5)

	// Angle at start: from posA toward control point
	angleA := -math.Atan2(float64(ctrl.Y-posA.ToF32().Y), float64(ctrl.X-posA.ToF32().X)) + math.Pi

	// Angle at end: from control point toward posB
	angleB := -math.Atan2(float64(posB.ToF32().Y-ctrl.Y), float64(posB.ToF32().X-ctrl.X))

	// Draw the arc
	DrawArc(ops, MoveAlongAngleGlob(posA, angleA+math.Pi, arrowSize*.5), MoveAlongAngleGlob(posB, angleB+math.Pi, arrowSize*.5), col, thickness, roundness, curvature)

	// Draw arrow at posA
	DrawArrow(ops, posA, angleA, arrowSize, col, windowSize)

	// Draw arrow at posB
	DrawArrow(ops, posB, angleB, arrowSize, col, windowSize)
}

func DrawArc(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness, roundness float32, curvature bool) {
	ctrl := GetCtrlPoint(posA.ToF32(), posB.ToF32(), roundness, curvature)

	var path clip.Path
	path.Begin(ops)
	path.MoveTo(posA.ToF32())
	path.QuadTo(ctrl, posB.ToF32())

	paint.FillShape(ops, col,
		clip.Stroke{
			Path:  path.End(),
			Width: thickness,
		}.Op(),
	)
}

func GetCtrlPoint(posA, posB f32.Point, roundness float32, curvature bool) f32.Point {
	mid := posA.Add(posB).Div(2)

	dx := float32(posB.X - posA.X)
	dy := float32(posB.Y - posA.Y)

	if curvature {
		return f32.Point{
			X: mid.X + dy*roundness,
			Y: mid.Y - dx*roundness,
		}
	} else {
		return f32.Point{
			X: mid.X - dy*roundness,
			Y: mid.Y + dx*roundness,
		}
	}
}

func DrawArrowLine(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness float32, windowSize GlobalDim) {
	angle := GetAngleGlob(posA, posB)
	arrowSize := float64(thickness * 5)

	DrawLine(ops, posA, MoveAlongAngleGlob(posB, angle+math.Pi, arrowSize*.5), col, thickness)
	DrawArrow(ops, posB, angle, arrowSize, col, windowSize)
}

func DrawLine(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness float32) {
	var path clip.Path
	path.Begin(ops)
	path.MoveTo(posA.ToF32())
	path.LineTo(posB.ToF32())

	paint.FillShape(ops, col,
		clip.Stroke{
			Path:  path.End(),
			Width: thickness,
		}.Op(),
	)
}

func DrawArrow(ops *op.Ops, basePos GlobalPos, angle float64, arrowSize float64, col color.NRGBA, windowSize GlobalDim) {
	var triangle clip.Path
	triangle.Begin(ops)
	triangle.MoveTo(MoveAlongAngleGlob(basePos, angle+math.Pi+math.Pi/7.0, arrowSize).ToF32())
	triangle.LineTo(MoveAlongAngleGlob(basePos, angle+math.Pi-math.Pi/7.0, arrowSize).ToF32())
	triangle.LineTo(basePos.ToF32())
	triangle.Close()

	defer clip.Outline{Path: triangle.End()}.Op().Push(ops).Pop()
	defer clip.Rect{Max: image.Pt(windowSize.W, windowSize.H)}.Push(ops).Pop()
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
