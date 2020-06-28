package neopixel

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/space"
)

// Line is a representation of a strait line of neopixels.
// Line implements effect.Device
type Line struct {
	// Row provides the implementation of effect.Allocater
	*Row

	// Start is the location of the first LED in the Line
	Start space.Vector

	// Theta is rotation about Z which the line continues at
	Theta float64
	// Phi is tilt from Z which the line continues at
	Phi float64
}

// NewLine creates a new Line
func NewLine(
	start space.Vector,
	theta float64,
	phi float64,
	length int,
) *Line {

	d := &Line{
		Start: start,
		Theta: theta,
		Phi:   phi,
	}

	singleLEDVector := space.NewVector(theta, phi, distPerLED)

	d.Row = NewRow(
		length,
		func() []*effect.Light {
			lights := make([]*effect.Light, length)
			for i := range lights {
				lights[i] = &effect.Light{
					Position: i,
					GetLocationFunc: func(position int) space.Vector {
						return start.Translate(singleLEDVector.Scale(float64(position)))
					},
				}
			}
			return lights
		},
	)

	return d
}
