package pdf

import (
	"image/color"
	"main/utils"
	"math"

	"github.com/jung-kurt/gofpdf"
)

func DrawRect(pdf *gofpdf.Fpdf, pos utils.LocalPos, dim utils.LocalDim, col color.NRGBA, thickness float32) {
	// Set fill color
	pdf.SetFillColor(int(col.R), int(col.G), int(col.B))

	// Draw filled rectangle
	pdf.Rect(float64(pos.X), float64(pos.Y), float64(dim.W), float64(dim.H), "F")

	// Draw outline if thickness > 0
	if thickness > 0 {
		pdf.SetLineWidth(float64(thickness))
		pdf.SetDrawColor(0, 0, 0) // Black outline
		pdf.Rect(float64(pos.X), float64(pos.Y), float64(dim.W), float64(dim.H), "D")
	}
}

func DrawEllipse(pdf *gofpdf.Fpdf, pos utils.LocalPos, dim utils.LocalDim, col color.NRGBA, thickness float32) {
	// Calculate center and radii
	cx := pos.X + dim.W/2
	cy := pos.Y + dim.H/2
	rx := dim.W / 2
	ry := dim.H / 2

	// Set fill color
	pdf.SetFillColor(int(col.R), int(col.G), int(col.B))

	// Draw filled ellipse
	pdf.Ellipse(float64(cx), float64(cy), float64(rx), float64(ry), 0, "F")

	// Draw outline
	pdf.SetLineWidth(float64(thickness))
	pdf.SetDrawColor(0, 0, 0) // Black outline
	pdf.Ellipse(float64(cx), float64(cy), float64(rx), float64(ry), 0, "D")
}

func DrawArrowLine(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32) {
	angle := utils.GetAngleLoc(posA, posB)
	arrowSize := thickness * 5

	// Draw line shortened at posB to accommodate arrow
	endPos := utils.MoveAlongAngleLoc(posB, angle+math.Pi, arrowSize*0.5)
	DrawLine(pdf, posA, endPos, col, thickness)

	// Draw arrow head at posB
	DrawArrowHead(pdf, posB, angle, arrowSize, col)
}

func DrawArrowCurve(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32, curvature float32) {
	// Calculate control point for quadratic bezier
	ctrl := utils.GetCtrlPoint(posA.ToF32(), posB.ToF32(), curvature)

	arrowSize := thickness * 5

	// Angle at start: from posA toward control point
	angleA := -math.Atan2(float64(ctrl.Y-posA.Y), float64(ctrl.X-posA.X)) + math.Pi

	// Angle at end: from control point toward posB
	angleB := -math.Atan2(float64(posB.Y-ctrl.Y), float64(posB.X-ctrl.X))

	// Draw the arc shortened at both ends
	startPos := utils.MoveAlongAngleLoc(posA, angleA+math.Pi, arrowSize*0.5)
	endPos := utils.MoveAlongAngleLoc(posB, angleB+math.Pi, arrowSize*0.5)
	DrawCurve(pdf, startPos, endPos, col, thickness, curvature)

	// Draw arrow heads
	DrawArrowHead(pdf, posA, angleA, arrowSize, col)
	DrawArrowHead(pdf, posB, angleB, arrowSize, col)
}

func DrawLine(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32) {
	pdf.SetLineWidth(float64(thickness))
	pdf.SetDrawColor(int(col.R), int(col.G), int(col.B))
	pdf.Line(float64(posA.X), float64(posA.Y), float64(posB.X), float64(posB.Y))
}

func DrawCurve(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32, curvature float32) {
	ctrl := utils.GetCtrlPoint(posA.ToF32(), posB.ToF32(), curvature)

	pdf.SetLineWidth(float64(thickness))
	pdf.SetDrawColor(int(col.R), int(col.G), int(col.B))

	// Draw quadratic bezier curve
	// gofpdf's CurveBezierCubic requires cubic bezier, so we convert quadratic to cubic
	// For quadratic P0, P1, P2: cubic control points are:
	// C1 = P0 + 2/3*(P1-P0)
	// C2 = P2 + 2/3*(P1-P2)
	c1x := posA.X + (2.0/3.0)*(ctrl.X-posA.X)
	c1y := posA.Y + (2.0/3.0)*(ctrl.Y-posA.Y)
	c2x := posB.X + (2.0/3.0)*(ctrl.X-posB.X)
	c2y := posB.Y + (2.0/3.0)*(ctrl.Y-posB.Y)

	pdf.CurveBezierCubic(float64(posA.X), float64(posA.Y), float64(c1x), float64(c1y), float64(c2x), float64(c2y), float64(posB.X), float64(posB.Y), "D")
}

func DrawArrowHead(pdf *gofpdf.Fpdf, basePos utils.LocalPos, angle float64, size float32, col color.NRGBA) {
	// Calculate triangle points
	p1 := utils.MoveAlongAngleLoc(basePos, angle+math.Pi+math.Pi/7.0, size)
	p2 := utils.MoveAlongAngleLoc(basePos, angle+math.Pi-math.Pi/7.0, size)

	// Set fill color
	pdf.SetFillColor(int(col.R), int(col.G), int(col.B))

	// Draw filled triangle
	pdf.Polygon([]gofpdf.PointType{
		{X: float64(p1.X), Y: float64(p1.Y)},
		{X: float64(p2.X), Y: float64(p2.Y)},
		{X: float64(basePos.X), Y: float64(basePos.Y)},
	}, "F")
}

func DrawArc(pdf *gofpdf.Fpdf, posA, posB, refPoint utils.LocalPos, radius float32, offsetAngle float64, col color.NRGBA, thickness float32) {
	circleCenter := utils.FindCircleCenter(posA.ToF32(), posB.ToF32(), refPoint.ToF32(), radius)
	angleA := utils.GetAngle(circleCenter, posA.ToF32())
	angleB := utils.GetAngle(circleCenter, posB.ToF32())

	offsetAngle *= 0.9 // shorten the offset angle a bit

	truncatedPosB := utils.MoveAlongAngle(circleCenter, utils.NormalizeAngle(angleB-offsetAngle), radius)

	angleDiff := math.Mod(angleB-angleA+math.Pi, 2*math.Pi) - math.Pi
	angle := utils.NormalizeAngle(angleDiff - 2*offsetAngle)

	pdf.SetLineWidth(float64(thickness))
	pdf.SetDrawColor(int(col.R), int(col.G), int(col.B))

	// Arc draws from angle1 to angle2 (in degrees) counter-clockwise
	startAngleDeg := math.Atan2(float64(truncatedPosB.Y-circleCenter.Y), float64(truncatedPosB.X-circleCenter.X)) * 180 / math.Pi
	sweepAngleDeg := angle * 180 / math.Pi

	pdf.SetLineCapStyle("round")
	pdf.Arc(float64(circleCenter.X), float64(circleCenter.Y), float64(radius), float64(radius), 0, startAngleDeg, startAngleDeg+sweepAngleDeg, "D")
}

func DrawArrowArc(pdf *gofpdf.Fpdf, posA, posB, refPoint utils.LocalPos, radius float32, col color.NRGBA, thickness float32) {
	circleCenter := utils.FindCircleCenter(posA.ToF32(), posB.ToF32(), refPoint.ToF32(), radius)
	arrowSize := float64(thickness * 5)
	offsetAngle := arrowSize / float64(radius)
	angleA := utils.GetAngle(circleCenter, posA.ToF32())
	angleB := utils.GetAngle(circleCenter, posB.ToF32())
	angleTangentA := utils.NormalizeAngle(angleA + offsetAngle - math.Pi/2)
	angleTangentB := utils.NormalizeAngle(angleB - offsetAngle + math.Pi/2)

	DrawArc(pdf, posA, posB, refPoint, radius, arrowSize/float64(radius), col, thickness)

	truncatedPosA := utils.MoveAlongAngle(circleCenter, utils.NormalizeAngle(angleA+offsetAngle), radius)
	truncatedPosB := utils.MoveAlongAngle(circleCenter, utils.NormalizeAngle(angleB-offsetAngle), radius)

	arrowPosA := utils.MoveAlongAngle(truncatedPosA, angleTangentA, float32(arrowSize))
	arrowPosB := utils.MoveAlongAngle(truncatedPosB, angleTangentB, float32(arrowSize))

	DrawArrowHead(pdf, utils.ToLocalPos(arrowPosA), angleTangentA, float32(arrowSize), col)
	DrawArrowHead(pdf, utils.ToLocalPos(arrowPosB), angleTangentB, float32(arrowSize), col)
}

func DrawText(pdf *gofpdf.Fpdf, pos utils.LocalPos, txt string, fontFamily string, bold bool, size, ppRatio float32) {
	styleStr := ""
	if bold {
		styleStr = "B"
	}

	// Set font
	pdf.SetFont(fontFamily, styleStr, float64(size*ppRatio)) // convert from Sp to PDF points

	// Set text color to black
	pdf.SetTextColor(0, 0, 0)

	// Position and draw text
	pdf.SetXY(float64(pos.X), float64(pos.Y))
	pdf.Cell(0, 0, txt)
}
