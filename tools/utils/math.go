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

func FindCircleCenter(p1, p2, sidePoint f32.Point, radius float32) f32.Point {
	// Midpoint of p1 and p2
	mx := (p1.X + p2.X) / 2
	my := (p1.Y + p2.Y) / 2

	// Distance between p1 and p2
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Check if circle is possible
	if dist > 2*radius {
		return f32.Point{} // Points too far apart
	}

	// Distance from midpoint to center along perpendicular bisector
	h := float32(math.Sqrt(float64(radius*radius - (dist/2)*(dist/2))))

	// Unit perpendicular vector
	px := -dy / dist
	py := dx / dist

	// Two possible centers
	c1 := f32.Point{X: mx + h*px, Y: my + h*py}
	c2 := f32.Point{X: mx - h*px, Y: my - h*py}

	// Choose center on same side as sidePoint
	// Use cross product to determine which side
	cross1 := (c1.X-p1.X)*(sidePoint.Y-p1.Y) - (c1.Y-p1.Y)*(sidePoint.X-p1.X)
	cross2 := (c2.X-p1.X)*(sidePoint.Y-p1.Y) - (c2.Y-p1.Y)*(sidePoint.X-p1.X)

	if cross1*cross2 > 0 {
		return c2
	}
	return c1
}

type Numeric interface {
	~int | ~float32 | ~float64
}

func RemapValue[T Numeric](val, iLow, iHigh, oLow, oHigh T) T {
	p := (val - iLow) / (iHigh - iLow)
	return p*(oHigh-oLow) + oLow
}

func AngleRectIntersection(angle float64, pos LocalPos, dim LocalDim) LocalPos {
	t := float32(math.Mod(angle+math.Pi/4, math.Pi/2))
	switch {
	case angle >= math.Pi/4 && angle < 3*math.Pi/4:
		return LocalPos{
			X: RemapValue(t, 0, math.Pi/2, pos.X+dim.W/2, pos.X-dim.W/2),
			Y: pos.Y - dim.H/2.0,
		}
	case angle >= 3*math.Pi/4 && angle < 5*math.Pi/4:
		return LocalPos{
			X: pos.X - dim.W/2.0,
			Y: RemapValue(t, 0, math.Pi/2, pos.Y-dim.H/2, pos.Y+dim.W/2),
		}
	case angle >= 5*math.Pi/4 && angle < 7*math.Pi/4:
		return LocalPos{
			X: RemapValue(t, 0, math.Pi/2, pos.X-dim.W/2, pos.X+dim.W/2),
			Y: pos.Y + dim.H/2.0,
		}
	case angle >= 7*math.Pi/4 || angle < math.Pi/4:
		return LocalPos{
			X: pos.X + dim.W/2.0,
			Y: RemapValue(t, 0, math.Pi/2, pos.Y+dim.H/2, pos.Y-dim.W/2),
		}
	default:
		return LocalPos{}
	}
}
