package utils

import "math"

func GetAngle(a, b LocalPos) float64 {
	rawAngle := -math.Atan2(float64(b.Y-a.Y), float64(b.X-a.X))
	return NormalizeAngle(rawAngle)
}

func NormalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 2 * math.Pi
	}
	for angle >= 2*math.Pi {
		angle -= 2 * math.Pi
	}
	return angle
}
