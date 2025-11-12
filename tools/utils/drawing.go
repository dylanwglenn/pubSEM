package utils

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/image/math/fixed"
)

type CoefficientDisplay int

const (
	NONE CoefficientDisplay = iota
	VALUE
	INTERVAL
	STAR
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
	if thickness > 0 {
		paint.FillShape(ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			clip.Stroke{
				Path:  rect.Path(),
				Width: thickness,
			}.Op(),
		)
	}
}

func DrawRoundedRect(ops *op.Ops, pos GlobalPos, dim GlobalDim, r int, col color.NRGBA, thickness float32) {
	rrect := clip.RRect{
		Rect: MakeRect(pos, dim),
		SE:   r,
		SW:   r,
		NE:   r,
		NW:   r,
	}

	defer rrect.Push(ops).Pop()

	// Draw fill
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)

	// Draw outline
	paint.FillShape(ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		clip.Stroke{
			Path:  rrect.Path(ops),
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

func DrawArrowArc(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness, curvature float32, windowSize GlobalDim) {
	// Calculate control point for tangent angles
	ctrl := GetCtrlPoint(posA.ToF32(), posB.ToF32(), curvature)

	arrowSize := float64(thickness * 5)

	// Angle at start: from posA toward control point
	angleA := -math.Atan2(float64(ctrl.Y-posA.ToF32().Y), float64(ctrl.X-posA.ToF32().X)) + math.Pi

	// Angle at end: from control point toward posB
	angleB := -math.Atan2(float64(posB.ToF32().Y-ctrl.Y), float64(posB.ToF32().X-ctrl.X))

	// Draw the arc
	DrawArc(ops, MoveAlongAngleGlob(posA, angleA+math.Pi, arrowSize*.5), MoveAlongAngleGlob(posB, angleB+math.Pi, arrowSize*.5), col, thickness, curvature)

	// Draw arrow at posA
	DrawArrowHead(ops, posA, angleA, arrowSize, col, windowSize)

	// Draw arrow at posB
	DrawArrowHead(ops, posB, angleB, arrowSize, col, windowSize)
}

func DrawArc(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness, curvature float32) {
	ctrl := GetCtrlPoint(posA.ToF32(), posB.ToF32(), curvature)

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

func GetCtrlPoint(posA, posB f32.Point, curvature float32) f32.Point {
	mid := posA.Add(posB).Div(2)

	dx := posB.X - posA.X
	dy := posB.Y - posA.Y

	return f32.Point{
		X: mid.X - dy*curvature,
		Y: mid.Y + dx*curvature,
	}
}

func DrawArrowLine(ops *op.Ops, posA, posB GlobalPos, col color.NRGBA, thickness float32, windowSize GlobalDim) {
	angle := GetAngleGlob(posA, posB)
	arrowSize := float64(thickness * 5)

	DrawLine(ops, posA, MoveAlongAngleGlob(posB, angle+math.Pi, arrowSize*.5), col, thickness)
	DrawArrowHead(ops, posB, angle, arrowSize, col, windowSize)
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

func DrawArrowHead(ops *op.Ops, basePos GlobalPos, angle float64, size float64, col color.NRGBA, windowSize GlobalDim) {
	var triangle clip.Path
	triangle.Begin(ops)
	triangle.MoveTo(MoveAlongAngleGlob(basePos, angle+math.Pi+math.Pi/7.0, size).ToF32())
	triangle.LineTo(MoveAlongAngleGlob(basePos, angle+math.Pi-math.Pi/7.0, size).ToF32())
	triangle.LineTo(basePos.ToF32())
	triangle.Close()

	defer clip.Outline{Path: triangle.End()}.Op().Push(ops).Pop()
	defer clip.Rect{Max: image.Pt(windowSize.W, windowSize.H)}.Push(ops).Pop()
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func DrawText(ops *op.Ops, gtx layout.Context, pos GlobalPos, txt string, style font.FontFace, size unit.Sp, scale float32) {
	defer op.Offset(pos.ToImagePnt()).Push(ops).Pop()

	// Apply scale transform
	defer op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(scale, scale))).Push(ops).Pop()

	// Create a label with the text
	label := material.Label(material.NewTheme(), size, txt)
	label.Font = style.Font

	// Draw the label
	label.Layout(gtx)
}

func GetTextWidth(txt string, style font.FontFace, size float32) float32 {
	shaper := text.NewShaper()
	params := text.Parameters{
		Font:    style.Font,
		PxPerEm: fixed.I(int(size)),
	}

	shaper.LayoutString(params, txt)

	var width fixed.Int26_6
	for {
		g, ok := shaper.NextGlyph()
		if !ok {
			break
		}
		width += g.Advance
	}

	spaceRunes := strings.Count(txt, " ")
	spaceWidth := 3                                             // this is a magic number. It was arrived at through trial and error
	return float32(width)/64.0 + float32(spaceRunes*spaceWidth) // Convert from fixed.Int26_6 to float32 and add space characters
}

func DrawEstimate(ops *op.Ops, gtx layout.Context, pos GlobalPos, fontStyle font.FontFace, fontSize float32, displayStyle CoefficientDisplay, est, pVal float64, ci [2]float64, precision int, scaleFactor float32, padding float32) (string, LocalDim) {
	// define the string to be printed
	estText, dim, textWidth := CalculateEstimate(fontStyle, fontSize, displayStyle, est, pVal, ci, precision, padding)

	DrawRect(ops, pos, dim.ToGlobal(scaleFactor), color.NRGBA{255, 255, 255, 255}, 0)

	// draw text
	textOffset := LocalDim{W: textWidth/2.0 - padding, H: fontSize / 1.5}
	DrawText(ops, gtx, pos.SubDim(textOffset.ToGlobal(scaleFactor)), estText, fontStyle, unit.Sp(fontSize-2), scaleFactor)

	return estText, dim
}

func CalculateEstimate(fontStyle font.FontFace, fontSize float32, displayStyle CoefficientDisplay, est, pVal float64, ci [2]float64, precision int, padding float32) (string, LocalDim, float32) {
	// define the string to be printed
	var estText string

	floatFmtStr := "%." + strconv.Itoa(precision) + "f"
	switch displayStyle {
	case VALUE:
		estText = fmt.Sprintf(floatFmtStr, est)
	case INTERVAL:
		estText = fmt.Sprintf(floatFmtStr+floatFmtStr, ci[0], ci[1])
	case STAR:
		var stars string
		switch {
		case pVal < .001:
			stars = "***"
		case pVal < .01:
			stars = "**"
		case pVal < .05:
			stars = "*"
		}
		estText = fmt.Sprintf(floatFmtStr+"%s", est, stars)
	default:
	}

	// draw the background rectangle
	textWidth := GetTextWidth(estText, fontStyle, fontSize)
	adjWidth := textWidth + padding*3.0
	height := fontSize * 1.5 // todo: find a better way to determine height of text
	return estText, LocalDim{W: adjWidth, H: height}, textWidth
}
