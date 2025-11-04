package utils

import (
	"math"
)

func SnapToGrid(pos LocalPos, gridSize float32) LocalPos {
	posX := float32(math.Round(float64(pos.X)/float64(gridSize))) * gridSize
	posY := float32(math.Round(float64(pos.Y)/float64(gridSize))) * gridSize

	return LocalPos{X: posX, Y: posY}
}
