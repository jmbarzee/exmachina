package neopixel

import (
	"math"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

// Ring is a representation of a ring of neopixels.
// Ring implements effect.Device
type Ring struct {
	// Row provides the implementation of effect.Allocater
	*Row

	// Radius is the distance from the Center to any LED
	Radius float64
	// Center is the point about which the Ring is centered
	Center space.Vector
	// Orientation is the rotation and tilt of the Ring
	Orientation space.Orientation
}

// NewRing creates a new Ring
func NewRing(
	radius float64,
	center space.Vector,
	orientation space.Orientation,
) *Ring {
	d := &Ring{
		Radius:      radius,
		Center:      center,
		Orientation: orientation,
	}
	length := int(radius * 2 * math.Pi / distPerLED)

	rotateTheta := space.NewRotationMatrixX(orientation.Theta)
	rotatePhi := space.NewRotationMatrixZ(orientation.Phi)
	orientationMatrix := rotateTheta.Mult(rotatePhi)
	radPerLED := distPerLED / radius

	d.Row = NewRow(
		length,
		func() []light.Light {
			lights := make([]light.Light, length)
			for i := range lights {

				localPhi := radPerLED * float64(i)

				// Location of LED if Ring was in XZ-Plane with first LED on the positive X axis
				sin, cos := math.Sincos(float64(localPhi))
				location := space.Vector{
					X: radius * cos,
					Y: radius * sin,
					Z: 0,
				}
				// Transform to match rotation and tilt of ring
				rotatedLocation := location.Transform(orientationMatrix)
				// Translate to be relative to origin and not center
				lightLocation := center.Translate(rotatedLocation)

				lightOrientation := orientation.Rotate(localPhi)

				lights[i] = &light.Basic{
					Position:     i,
					NumPositions: length,
					Location:     lightLocation,
					Orientation:  lightOrientation,
				}
			}
			return lights
		},
	)
	return d
}
