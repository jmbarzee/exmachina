package color

func max(a, b, c float32) float32 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}

func min(a, b, c float32) float32 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// Average does what it says it does.
func Average(a, b float32) float32 {
	return (a + b) / 2
}

// AverageWeighted is like average but with weighting
// weight should range from 0 to 1.
// 0 will favor the first number and 1 will favor the second number.
func AverageWeighted(a, b, weight float32) float32 {
	first := a * (1 - weight)
	second := b * weight
	return (first + second)
}
