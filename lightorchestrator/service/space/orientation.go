package space

import "math"

// Orientation represents the direction of spherical coordinates
type Orientation struct {
	// Phi is rotation about Z
	Phi float32
	// Theta is tilt from Z
	Theta float32
}

// NewOrientation creates a new Orientation from a rotation and tilt
func NewOrientation(phi, theta float32) Orientation {
	o := Orientation{}
	o = o.Rotate(phi)
	o = o.Tilt(theta)
	return o
}

// Rotate will adjust the rotation about Z by phi
func (o Orientation) Rotate(phi float32) Orientation {
	wrappedPhi := o.Phi + phi
	o.Phi = float32(math.Mod(float64(wrappedPhi), math.Pi*2))
	return o
}

// Tilt will adjust the tilt from Z by theta
func (o Orientation) Tilt(theta float32) Orientation {
	wrappedTheta := o.Theta + theta
	newTheta := float32(math.Mod(float64(wrappedTheta), math.Pi*2))

	// Check if tilt is negative
	if newTheta < 0 {
		newTheta = (math.Pi * 2) + newTheta
	}

	// Check if tilt is beyond range of [0, pi]
	if newTheta > math.Pi {
		o.Theta = (math.Pi * 2) - newTheta
		return o.Rotate(math.Pi)
	}

	o.Theta = newTheta
	return o
}
