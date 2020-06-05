package helpers

import "math"

// Sinf32 retuns sine of float32 argument in radians
func Sinf32(rad float32) float32 {
	return float32(math.Sin(float64(rad)))
}

// Cosf32 retuns cossine of float32 argument in radians
func Cosf32(rad float32) float32 {
	return float32(math.Cos(float64(rad)))
}
