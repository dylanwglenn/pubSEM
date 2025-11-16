package pdf

import "C"
import (
	"embed"
	"image/color"
	"main/model"
	"main/utils"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

const (
	docPadding = 15
	ppRatio    = .75  // pixel-to-point conversion
	textAdj    = 3.42 // this is seemingly random. Arrived at through trial and error
)

//go:embed "gofpdf_fonts"
var pdfFontFS embed.FS

func ExportModel(m *model.Model, filePath string) {
	rect, localDim := GetModelSize(m)

	pageWidth := localDim.W*ppRatio + 2*docPadding
	pageHeight := localDim.H*ppRatio + 2*docPadding

	fontDir := createTempFontDir()
	defer os.RemoveAll(fontDir)

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "pt",
		Size:       gofpdf.SizeType{Wd: float64(pageWidth), Ht: float64(pageHeight)},
		FontDirStr: fontDir,
	})
	pdf.SetAutoPageBreak(false, 0) // required to avoid automatic page breaks at different text sizes
	utils.LoadPdfFonts(pdf)
	pdf.AddPage()

	// Offset to translate model coordinates to page coordinates
	offsetX := docPadding - rect[0].X
	offsetY := docPadding - rect[0].Y

	for _, n := range m.Nodes {
		if !n.Visible {
			continue
		}

		// Convert center position to top-left corner
		adjPos := utils.LocalPos{
			X: (n.Pos.X - n.Dim.W/2 + offsetX) * ppRatio,
			Y: (n.Pos.Y - n.Dim.H/2 + offsetY) * ppRatio,
		}

		adjDim := n.Dim.Mul(ppRatio)

		switch n.Class {
		case model.OBSERVED:
			DrawRect(pdf, adjPos, adjDim, n.Col, n.Thickness*ppRatio*.5) // .5 adjusts thickness from Gio to gofpdf
		case model.LATENT:
			DrawEllipse(pdf, adjPos, adjDim, n.Col, n.Thickness*ppRatio*.5)
		case model.INTERCEPT:
			// todo: handle intercepts
		}

		textPos := utils.LocalPos{
			X: adjPos.X - textAdj + n.Padding*ppRatio,
			Y: adjPos.Y + adjDim.H/2,
		}

		DrawText(pdf, textPos, n.Text, m.Font.Family, n.Bold, m.Font.Size, ppRatio)
	}

	for _, c := range m.Connections {
		// adjust connection points to PDS coords
		originPos := utils.LocalPos{
			X: (c.OriginPos.X + offsetX) * ppRatio,
			Y: (c.OriginPos.Y + offsetY) * ppRatio,
		}
		destPos := utils.LocalPos{
			X: (c.DestinationPos.X + offsetX) * ppRatio,
			Y: (c.DestinationPos.Y + offsetY) * ppRatio,
		}

		switch c.Type {
		case model.STRAIGHT:
			DrawArrowLine(pdf, originPos, destPos, c.Col, c.Thickness*ppRatio)
		case model.CURVED:
			DrawArrowCurve(pdf, originPos, destPos, c.Col, c.Thickness*ppRatio, c.Curvature)
		case model.CIRCULAR:
			refPos := utils.LocalPos{
				X: (c.RefPos.X + offsetX) * ppRatio,
				Y: (c.RefPos.Y + offsetY) * ppRatio,
			}
			var radius float32 = 20
			DrawArrowArc(pdf, originPos, destPos, refPos, radius, c.Col, c.Thickness*ppRatio)
		}
	}

	// draw estimate labels after ALL of the connections to ensure proper layering
	for _, c := range m.Connections {
		textWidth := utils.GetTextWidth(c.EstText, m.Font.Face, (m.Font.Size-2)*ppRatio) + (c.EstPadding * ppRatio)
		textPos := utils.LocalPos{
			X: (c.EstPos.X+offsetX)*ppRatio - textWidth/2 - textAdj,
			Y: (c.EstPos.Y+offsetY)*ppRatio - ppRatio/2, // assuming that ppRatio is text height
		}

		rectPos := utils.LocalPos{
			X: (c.EstPos.X-c.EstDim.W/2+offsetX)*ppRatio + (c.EstPadding * ppRatio),
			Y: (c.EstPos.Y - c.EstDim.H/2 + offsetY) * ppRatio,
		}

		rectDim := c.EstDim.Mul(ppRatio)

		DrawRect(pdf, rectPos, rectDim, color.NRGBA{255, 255, 255, 255}, 0)
		DrawText(pdf, textPos, c.EstText, m.Font.Family, false, m.Font.Size-2, ppRatio)
	}

	// export
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		panic(err)
	}
}

func GetModelSize(m *model.Model) (rect [2]utils.LocalPos, dim utils.LocalDim) {
	// first LocalPos in rect is the NW corner, second is the SE corner
	// initialize rect as an existing position to ensure resultant rect is directly against the shapes
	rect = [2]utils.LocalPos{m.Nodes[0].Pos, m.Nodes[0].Pos}

	for _, n := range m.Nodes {
		minX := n.Pos.X - n.Dim.W/2
		maxX := n.Pos.X + n.Dim.W/2
		minY := n.Pos.Y - n.Dim.H/2
		maxY := n.Pos.Y + n.Dim.H/2

		// handle x coords
		if minX < rect[0].X {
			rect[0].X = minX
		}
		if maxX > rect[1].X {
			rect[1].X = maxX
		}

		// handle y coords
		// (remember that lower ys are visually higher)
		if minY < rect[0].Y {
			rect[0].Y = minY
		}
		if maxY > rect[1].Y {
			rect[1].Y = maxY
		}
	}

	for _, c := range m.Connections {
		minX := c.EstPos.X - c.EstDim.W/2
		maxX := c.EstPos.X + c.EstDim.W/2
		minY := c.EstPos.Y - c.EstDim.H/2
		maxY := c.EstPos.Y + c.EstDim.H/2

		// handle x coords
		if minX < rect[0].X {
			rect[0].X = minX
		}
		if maxX > rect[1].X {
			rect[1].X = maxX
		}

		// handle y coords
		// (remember that lower ys are visually higher)
		if minY < rect[0].Y {
			rect[0].Y = minY
		}
		if maxY > rect[1].Y {
			rect[1].Y = maxY
		}
	}

	dim = utils.LocalDim{
		W: utils.Abs32(rect[1].X - rect[0].X),
		H: utils.Abs32(rect[1].Y - rect[0].Y),
	}

	return
}

func createTempFontDir() string {
	tempDir, err := os.MkdirTemp("", "gofpdf_fonts_*")
	if err != nil {
		panic(err)
	}

	fontFiles, err := pdfFontFS.ReadDir("gofpdf_fonts")
	if err != nil {
		panic(err)
	}
	for _, f := range fontFiles {
		if f.IsDir() {
			continue
		}

		data, err := pdfFontFS.ReadFile(filepath.Join("gofpdf_fonts", f.Name()))
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(filepath.Join(tempDir, f.Name()), data, 0644)
		if err != nil {
			panic(err)
		}
	}

	return tempDir
}
