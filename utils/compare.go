package utils

import "math"

const float64EqualityThreshold = 1e-9

// AlmostEqual 比较float64
func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
