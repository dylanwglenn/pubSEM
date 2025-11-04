package utils

import "image"

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
