package utils

import (
	"image"

	"gioui.org/f32"
)

func WithinRect(pos image.Point, rect image.Rectangle) bool {
	return (pos.Y > rect.Min.Y && pos.Y < rect.Max.Y) && (pos.X > rect.Min.X && pos.X < rect.Max.X)
}

func WithinEllipse(pos image.Point, rect image.Rectangle) bool {
	// Center of the ellipse
	cx := float64(rect.Min.X+rect.Max.X) / 2.0
	cy := float64(rect.Min.Y+rect.Max.Y) / 2.0

	// Semi-axes (radii)
	rx := float64(rect.Max.X-rect.Min.X) / 2.0
	ry := float64(rect.Max.Y-rect.Min.Y) / 2.0

	// Ellipse equation: ((x-cx)/rx)^2 + ((y-cy)/ry)^2 <= 1
	dx := (float64(pos.X) - cx) / rx
	dy := (float64(pos.Y) - cy) / ry

	return dx*dx+dy*dy <= 1.0
}

func WithinLine(pos image.Point, a, b GlobalPos, tolerance float32) bool {
	p := f32.Point{X: float32(pos.X), Y: float32(pos.Y)}
	// Vector from a to b
	ab := b.ToF32().Sub(a.ToF32())
	// Vector from a to p
	ap := p.Sub(a.ToF32())

	// Project p onto line segment
	abLenSq := ab.X*ab.X + ab.Y*ab.Y
	if abLenSq == 0 {
		// a and b are the same point
		dx := pos.X - a.X
		dy := pos.Y - a.Y
		return float32(dx*dx+dy*dy) <= tolerance*tolerance
	}

	t := (ap.X*ab.X + ap.Y*ab.Y) / abLenSq

	// Clamp t to [0, 1] to stay on segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Closest point on segment
	closest := f32.Point{
		X: float32(a.X) + t*float32(ab.X),
		Y: float32(a.Y) + t*float32(ab.Y),
	}

	// Distance from p to nearest point
	dx := p.X - closest.X
	dy := p.Y - closest.Y
	distSq := dx*dx + dy*dy

	return distSq <= tolerance*tolerance
}

func WithinArc(pos image.Point, a, b GlobalPos, roundness, tolerance float32, curvature bool, samples int) bool {
	ctrl := GetCtrlPoint(a.ToF32(), b.ToF32(), roundness, curvature)

	// Sample the curve and check distance to each segment
	// More samples = better accuracy, but slower
	for i := 0; i < samples; i++ {
		t1 := float32(i) / float32(samples)
		t2 := float32(i+1) / float32(samples)

		p1 := evalQuadraticBezier(a.ToF32(), ctrl, b.ToF32(), t1)
		p2 := evalQuadraticBezier(a.ToF32(), ctrl, b.ToF32(), t2)

		if WithinLine(pos, ToGlobalPos(p1.Round()), ToGlobalPos(p2.Round()), tolerance) {
			return true
		}
	}

	return false
}

func evalQuadraticBezier(p0, p1, p2 f32.Point, t float32) f32.Point {
	s := 1 - t
	return f32.Point{
		X: s*s*p0.X + 2*s*t*p1.X + t*t*p2.X,
		Y: s*s*p0.Y + 2*s*t*p1.Y + t*t*p2.Y,
	}
}
