package utils

import "math"

func Abs32(x float32) float32 {
	return float32(math.Abs(float64(x)))
}
