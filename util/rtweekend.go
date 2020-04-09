package util

import (
	"math/rand"
)

// Constants
const pi float64 = 3.1415926535897932385

// Utility Functions
func DegreesToRadians(degrees float64) float64 {
	return degrees * pi / 180
}

func FFMin(a float64, b float64) float64 { if a <= b { return a } else { return b } }
func FFMax(a float64, b float64) float64 { if a >= b { return a } else { return b } }

func Clamp(x float64, min float64, max float64) float64 {
	if x < min { return min }
	if x > max { return max }
	return x
}

func Random() float64 {
	// Returns a random float64 in [0, 1).
	return rand.Float64()
}

func RandomBetween(min float64, max float64) float64 {
	// Returns a random float64 in [min, max).
	return min + (max-min)*Random()
}
