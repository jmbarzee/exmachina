package neopixel

import (
	"math"

	"github.com/jmbarzee/dominion/services/lightorchestrator/service/shared"
	"github.com/jmbarzee/dominion/services/lightorchestrator/service/space"
)

// Ring is a representation of a ring of neopixels.
// Ring implements effect.Device
type Ring struct {
	// Row provides the implementation of effect.Allocater
	*Row

	// Center is the point about which the Ring is centered
	Center space.Vector
	// Radius is the distance from the Center to any LED
	Radius float32

	// Theta is the rotation about Z of the first LED
	Theta float32
	// Phi is the tilt from Z of the first LED
	Phi float32
}

// NewRing creates a new Ring
func NewRing(
	center space.Vector,
	radius float32,
	theta float32,
	phi float32,
) *Ring {
	d := &Ring{
		Center: center,
		Radius: radius,

		Theta: theta,
		Phi:   phi,
	}
	length := int(radius * 2 * math.Pi / distPerLED)

	transRotatePhi := space.NewRotationMatrixY(math.Pi/2.0 - phi)
	transRotateTheta := space.NewRotationMatrixZ(theta)
	transRotate := transRotateTheta.Mult(transRotatePhi)
	radPerLED := distPerLED / radius

	d.Row = NewRow(
		length,
		func() []shared.Light {
			lights := make([]shared.Light, length)
			for i := range lights {
				lights[i] = shared.Light{
					Position: i,
					GetLocationFunc: func(position int) space.Vector {

						radToLED := radPerLED * float32(position)

						// Location of LED if Ring was in XY-Plane with first LED on the positive X axis
						location := space.Vector{
							X: radius * float32(math.Cos(float64(radToLED))),
							Y: radius * float32(math.Sin(float64(radToLED))),
							Z: 0,
						}
						// Transform to match rotation and tilt of ring
						transLocation := location.Transform(transRotate)
						// Translate to be relative to origin and not center
						return center.Translate(transLocation)
					},
				}
			}
			return lights
		},
	)
	return d
}
