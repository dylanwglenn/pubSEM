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

func DrawArrowArc(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32, curvature float32) {
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
	DrawArc(pdf, startPos, endPos, col, thickness, curvature)

	// Draw arrow heads
	DrawArrowHead(pdf, posA, angleA, arrowSize, col)
	DrawArrowHead(pdf, posB, angleB, arrowSize, col)
}

func DrawLine(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32) {
	pdf.SetLineWidth(float64(thickness))
	pdf.SetDrawColor(int(col.R), int(col.G), int(col.B))
	pdf.Line(float64(posA.X), float64(posA.Y), float64(posB.X), float64(posB.Y))
}

func DrawArc(pdf *gofpdf.Fpdf, posA, posB utils.LocalPos, col color.NRGBA, thickness float32, curvature float32) {
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

func DrawText(pdf *gofpdf.Fpdf, pos utils.LocalPos, txt string, fontFamily string, bold bool, size float32) {
	styleStr := ""
	if bold {
		styleStr = "B"
	}
	// Set font
	pdf.SetFont(fontFamily, styleStr, float64(size*0.75)) // 0.75 is approximate conversion from Sp to PDF points

	// Set text color to black (adjust if needed)
	pdf.SetTextColor(0, 0, 0)

	// Position and draw text
	pdf.SetXY(float64(pos.X), float64(pos.Y))
	pdf.Cell(0, 0, txt)
}
