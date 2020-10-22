package color

import "math"

func max(a, b, c float64) float64 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}

func min(a, b, c float64) float64 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}
func modOne(val float64) float64 {
	return math.Mod(val, 1.0)
}

// Average does what it says it does.
func Average(a, b float64) float64 {
	return (a + b) / 2
}

// AverageWeighted is like average but with weighting
// weight should range from 0 to 1.
// 0 will favor the first number and 1 will favor the second number.
func AverageWeighted(a, b, weight float64) float64 {
	first := a * (1 - weight)
	second := b * weight
	return (first + second)
}
