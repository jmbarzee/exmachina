package space

import "math"

// Matrix is a transformational matrix for 3D space (3 x 3)
type Matrix [][]float32

// NewRotationMatrixX produces a matrix which will rotate about X
func NewRotationMatrixX(theta float32) Matrix {
	sin64, cos64 := math.Sincos(float64(theta))
	sin := float32(sin64)
	cos := float32(cos64)
	return Matrix{
		{1, 0, 0},
		{0, cos, -sin},
		{0, sin, cos},
	}
}

// NewRotationMatrixY produces a matrix which will rotate about X
func NewRotationMatrixY(theta float32) Matrix {
	sin64, cos64 := math.Sincos(float64(theta))
	sin := float32(sin64)
	cos := float32(cos64)
	return Matrix{
		{cos, 0, sin},
		{0, 1, 0},
		{-sin, 0, -cos},
	}
}

// NewRotationMatrixZ produces a matrix which will rotate about X
func NewRotationMatrixZ(theta float32) Matrix {
	sin64, cos64 := math.Sincos(float64(theta))
	sin := float32(sin64)
	cos := float32(cos64)
	return Matrix{
		{cos, -sin, 0},
		{sin, cos, 0},
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
