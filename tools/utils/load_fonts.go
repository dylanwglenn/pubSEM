package utils

import (
	_ "embed"

	"gioui.org/font/opentype"
	"gioui.org/text"
	"github.com/jung-kurt/gofpdf"
)

//go:embed fonts/Noto_Sans/static/NotoSans-Medium.ttf
var sansNormalData []byte

//go:embed fonts/Noto_Sans/static/NotoSans-Bold.ttf
var sansBoldData []byte

//go:embed fonts/Noto_Serif/static/NotoSerif-Medium.ttf
var serifNormalData []byte

//go:embed fonts/Noto_Serif/static/NotoSerif-Bold.ttf
var serifBoldData []byte

func LoadSansFontFace() []*text.FontFace {
	normal, _ := opentype.Parse(sansNormalData)
	bold, _ := opentype.Parse(sansBoldData)

	normalFontFace := text.FontFace{Font: normal.Font(), Face: normal}
	boldFontFace := text.FontFace{Font: bold.Font(), Face: bold}
	return []*text.FontFace{&normalFontFace, &boldFontFace}
}

func LoadSerifFontFace() []text.FontFace {
	normal, _ := opentype.Parse(serifNormalData)
	bold, _ := opentype.Parse(serifBoldData)

	normalFontFace := text.FontFace{Font: normal.Font(), Face: normal}
	boldFontFace := text.FontFace{Font: bold.Font(), Face: bold}
	return []text.FontFace{normalFontFace, boldFontFace}
}

func LoadPdfFonts(pdf *gofpdf.Fpdf) {
	// Add sans regular
	pdf.AddFont("sans", "", "NotoSans-Regular.json")
	// Add sans bold
	pdf.AddFont("sans", "B", "NotoSans-Bold.json")

	// Add serif regular
	pdf.AddFont("serif", "", "NotoSerif-Regular.json")
	// Add serif bold
	pdf.AddFont("serif", "B", "NotoSerif-Bold.json")
}
