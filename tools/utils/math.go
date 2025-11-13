package utils

import (
	"math"

	"gioui.org/f32"
)

func Abs32(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func ProjectOntoLine(a, b, p f32.Point) (LocalPos, float32) {
	// Vector from a to b
	ab := b.Sub(a)
	// Vector from a to p
	ap := p.Sub(a)

	// Project p onto line segment
	abLenSq := ab.X*ab.X + ab.Y*ab.Y

	t := (ap.X*ab.X + ap.Y*ab.Y) / abLenSq

	// Clamp t to [0, 1] to stay on segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Closest point on segment
	return LocalPos{
		X: a.X + t*ab.X,
		Y: a.Y + t*ab.Y,
	}, t
}

func MoveAlongBezier(a, b, ctrl f32.Point, t float32) LocalPos {
	// Use quadratic Bézier: (1-t)²P₀ + 2(1-t)tP₁ + t²P₂
	oneMinusT := 1.0 - t

	x := oneMinusT*oneMinusT*a.X + 2*oneMinusT*t*ctrl.X + t*t*b.X
	y := oneMinusT*oneMinusT*a.Y + 2*oneMinusT*t*ctrl.Y + t*t*b.Y

	return LocalPos{X: x, Y: y}
}

func UnitVector(a, b LocalPos) LocalPos {
	return LocalPos{X: b.X - a.X, Y: b.Y - a.Y}
}
