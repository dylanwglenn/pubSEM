package utils

import (
	"image"
	"math"

	"gioui.org/f32"
)

type LocalDim struct {
	W, H float32
}
type LocalPos struct {
	X, Y float32
}
type GlobalDim struct {
	W, H int
}
type GlobalPos struct {
	X, Y int
}

func DistLoc(a, b LocalPos) float32 {
	return float32(math.Sqrt(math.Pow(float64(b.X-a.X), 2) + math.Pow(float64(b.Y-a.Y), 2)))
}

// LocalDim methods

func (lDim LocalDim) Round() GlobalDim {
	wRound := int(math.Round(float64(lDim.W)))
	hRound := int(math.Round(float64(lDim.H)))
	return GlobalDim{W: wRound, H: hRound}
}

func (lDim LocalDim) Add(x LocalDim) LocalDim {
	return LocalDim{W: lDim.W + x.W, H: lDim.H + x.H}
}

func (lDim LocalDim) Sub(x LocalDim) LocalDim {
	return LocalDim{W: lDim.W - x.W, H: lDim.H - x.H}
}

func (lDim LocalDim) Mul(f float32) LocalDim {
	return LocalDim{W: lDim.W * f, H: lDim.H * f}
}

func (lDim LocalDim) Div(f float32) LocalDim {
	return LocalDim{W: lDim.W / f, H: lDim.H / f}
}

func (lDim LocalDim) ToGlobal(scaleFactor float32) GlobalDim {
	return lDim.Mul(scaleFactor).Round()
}

// LocalPos methods

func (lPos LocalPos) Round() GlobalPos {
	wRound := int(math.Round(float64(lPos.X)))
	hRound := int(math.Round(float64(lPos.Y)))
	return GlobalPos{X: wRound, Y: hRound}
}

func (lPos LocalPos) Add(x LocalPos) LocalPos {
	return LocalPos{X: lPos.X + x.X, Y: lPos.Y + x.Y}
}

func (lPos LocalPos) AddDim(x LocalDim) LocalPos {
	return LocalPos{X: lPos.X + x.W, Y: lPos.Y + x.H}
}

func (lPos LocalPos) Sub(x LocalPos) LocalPos {
	return LocalPos{X: lPos.X - x.X, Y: lPos.Y - x.Y}
}

func (lPos LocalPos) SubDim(x LocalDim) LocalPos {
	return LocalPos{X: lPos.X - x.W, Y: lPos.Y - x.H}
}

func (lPos LocalPos) Mul(f float32) LocalPos {
	return LocalPos{X: lPos.X * f, Y: lPos.Y * f}
}

func (lPos LocalPos) Div(f float32) LocalPos {
	return LocalPos{X: lPos.X / f, Y: lPos.Y / f}
}

func (lPos LocalPos) ToGlobal(scaleFactor float32, viewportCenter LocalPos, windowSize GlobalDim) GlobalPos {
	return lPos.Add(viewportCenter).Mul(scaleFactor).Round().AddDim(windowSize.Div(2))
}

func (lPos LocalPos) ToF32() f32.Point {
	return f32.Point{X: lPos.X, Y: lPos.Y}
}

func ToLocalPos(pt f32.Point) LocalPos {
	return LocalPos{X: pt.X, Y: pt.Y}
}

// GlobalDim methods

func (gDim GlobalDim) Add(x GlobalDim) GlobalDim {
	return GlobalDim{W: gDim.W + x.W, H: gDim.H + x.H}
}

func (gDim GlobalDim) Sub(x GlobalDim) GlobalDim {
	return GlobalDim{W: gDim.W - x.W, H: gDim.H - x.H}
}

func (gDim GlobalDim) Mul(i int) GlobalDim {
	return GlobalDim{W: gDim.W * i, H: gDim.H * i}
}

func (gDim GlobalDim) Div(i int) GlobalDim {
	return GlobalDim{W: gDim.W / i, H: gDim.H / i}
}

// GlobalPos methods

func (gPos GlobalPos) Add(x GlobalPos) GlobalPos {
	return GlobalPos{X: gPos.X + x.X, Y: gPos.Y + x.Y}
}

func (gPos GlobalPos) AddDim(x GlobalDim) GlobalPos {
	return GlobalPos{X: gPos.X + x.W, Y: gPos.Y + x.H}
}

func (gPos GlobalPos) SubDim(x GlobalDim) GlobalPos {
	return GlobalPos{X: gPos.X - x.W, Y: gPos.Y - x.H}
}

func (gPos GlobalPos) Sub(x GlobalPos) GlobalPos {
	return GlobalPos{X: gPos.X - x.X, Y: gPos.Y - x.Y}
}

func (gPos GlobalPos) Mul(i int) GlobalPos {
	return GlobalPos{X: gPos.X * i, Y: gPos.Y * i}
}

func (gPos GlobalPos) Div(i int) GlobalPos {
	return GlobalPos{X: gPos.X / i, Y: gPos.Y / i}
}

func (gPos GlobalPos) ToF32() f32.Point {
	return f32.Point{X: float32(gPos.X), Y: float32(gPos.Y)}
}

func (gPos GlobalPos) ToImagePnt() image.Point {
	return image.Point{X: gPos.X, Y: gPos.Y}
}

func ToGlobalPos(pt image.Point) GlobalPos {
	return GlobalPos{X: pt.X, Y: pt.Y}
}
