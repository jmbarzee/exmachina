package space

import "math"

// Matrix is a transformational matrix for 3D space (3 x 3)
type Matrix [][]float64

// NewRotationMatrixX produces a matrix which will rotate about X
func NewRotationMatrixX(theta float64) Matrix {
	return Matrix{
		{1, 0, 0},
		{0, math.Cos(theta), -math.Sin(theta)},
		{0, math.Sin(theta), math.Cos(theta)},
	}
}

// NewRotationMatrixY produces a matrix which will rotate about X
func NewRotationMatrixY(theta float64) Matrix {
	return Matrix{
		{math.Cos(theta), 0, math.Sin(theta)},
		{0, 1, 0},
		{-math.Sin(theta), 0, -math.Cos(theta)},
	}
}

// NewRotationMatrixZ produces a matrix which will rotate about X
func NewRotationMatrixZ(theta float64) Matrix {
	return Matrix{
		{math.Cos(theta), -math.Sin(theta), 0},
		{math.Sin(theta), math.Cos(theta), 0},
		{0, 0, 1},
	}
}

// Mult will return the result of m * n
func (m Matrix) Mult(n Matrix) Matrix {
	r := Matrix{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	for rowM := range m {
		for colN := 0; colN < 3; colN++ {
			a := m[rowM][0] * n[0][colN]
			b := m[rowM][1] * n[1][colN]
			c := m[rowM][2] * n[2][colN]
			r[rowM][colN] = a + b + c
		}
	}
	return r
}
