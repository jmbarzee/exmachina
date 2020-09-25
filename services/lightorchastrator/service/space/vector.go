package space

import "math"

// Vector is a 3D coordinate (also known as Point)
type Vector struct {
	X, Y, Z float32
}

// NewVector produces a new Vector from spherical coordinates
func NewVector(direction Orientation, radius float32) Vector {
	sinTheta64, cosTheta64 := math.Sincos(float64(direction.Theta))
	sinTheta := float32(sinTheta64)
	cosTheta := float32(cosTheta64)
	sinPhi64, cosPhi64 := math.Sincos(float64(direction.Phi))
	sinPhi := float32(sinPhi64)
	cosPhi := float32(cosPhi64)
	return Vector{
		X: radius * sinTheta * cosPhi,
		Y: radius * sinTheta * sinPhi,
		Z: radius * cosTheta,
	}
}

// Translate shifts a vector by another vector (addition)
func (v Vector) Translate(q Vector) Vector {
	return Vector{
		X: v.X + q.X,
		Y: v.Y + q.Y,
		Z: v.Z + q.Z,
	}
}

// Scale multivlies a vector by a given scale
func (v Vector) Scale(i float32) Vector {
	return Vector{
		X: v.X * i,
		Y: v.Y * i,
		Z: v.Z * i,
	}
}

func (v Vector) Transform(m Matrix) Vector {
	return Vector{
		X: (v.X * m[0][0]) + (v.Y * m[0][1]) + (v.Z * m[0][2]),
		Y: (v.X * m[1][0]) + (v.Y * m[1][1]) + (v.Z * m[1][2]),
		Z: (v.X * m[2][0]) + (v.Y * m[2][1]) + (v.Z * m[2][2]),
	}
}
