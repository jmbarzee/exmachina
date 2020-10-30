package space

import "math"

// Matrix is a transformational matrix for 3D space (3 x 3)
type Matrix [][]float64

// newRotationMatrixX produces a matrix which will rotate about X
func newRotationMatrixX(theta float64) Matrix {
	sin, cos := math.Sincos(theta)
	return Matrix{
		{1, 0, 0, 0},
		{0, cos, -sin, 0},
		{0, sin, cos, 0},
		{0, 0, 0, 1},
	}
}

// newRotationMatrixY produces a matrix which will rotate about Y
func newRotationMatrixY(theta float64) Matrix {
	sin, cos := math.Sincos(theta)
	return Matrix{
		{cos, 0, sin, 0},
		{0, 1, 0, 0},
		{-sin, 0, -cos, 0},
		{0, 0, 0, 1},
	}
}

// newRotationMatrixZ produces a matrix which will rotate about Z
func newRotationMatrixZ(theta float64) Matrix {
	sin, cos := math.Sincos(theta)
	return Matrix{
		{cos, -sin, 0, 0},
		{sin, cos, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Mult will return the result of m * n
func (m Matrix) Mult(n Matrix) Matrix {
	r := Matrix{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 1},
	}
	for rowM := range m {
		for colN := 0; colN < 4; colN++ {
			a := m[rowM][0] * n[0][colN]
			b := m[rowM][1] * n[1][colN]
			c := m[rowM][2] * n[2][colN]
			d := m[rowM][3] * n[3][colN]
			r[rowM][colN] = a + b + c + d
		}
	}
	return r
}
