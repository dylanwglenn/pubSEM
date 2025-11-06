package utils

import "math"

func GetAngleLoc(a, b LocalPos) float64 {
	rawAngle := -math.Atan2(float64(b.Y-a.Y), float64(b.X-a.X))
	return NormalizeAngle(rawAngle)
}

func GetAngleGlob(a, b GlobalPos) float64 {
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

func MoveAlongAngleGlob(pos GlobalPos, angle float64, dist float64) GlobalPos {
	newX := pos.X + int(dist*math.Cos(angle))
	newY := pos.Y - int(dist*math.Sin(angle))

	return GlobalPos{newX, newY}
}

func MoveAlongAngleLoc(pos LocalPos, angle float64, dist float32) LocalPos {
	newX := pos.X + dist*float32(math.Cos(angle))
	newY := pos.Y - dist*float32(math.Sin(angle))

	return LocalPos{newX, newY}
}
