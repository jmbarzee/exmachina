package neopixel

import (
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/shared"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/space"
)

// Line is a representation of a strait line of neopixels.
// Line implements effect.Device
type Line struct {
	// Row provides the implementation of effect.Allocater
	*Row

	// Start is the location of the first LED in the Line
	Start space.Vector
	// Direction is direction which all LEDs are from the first
	Direction space.Orientation
}

// NewLine creates a new Line
func NewLine(start space.Vector, direction space.Orientation, length int) *Line {

	d := &Line{
		Start:     start,
		Direction: direction,
	}

	singleLEDVector := space.NewVector(direction, distPerLED)

	d.Row = NewRow(
		length,
		func() []shared.Light {
			lights := make([]shared.Light, length)
			for i := range lights {
				lights[i] = shared.Light{
					Position: i,
					GetLocationFunc: func(position int) space.Vector {
						return start.Translate(singleLEDVector.Scale(float32(position)))
					},
				}
			}
			return lights
		},
	)

	return d
}

// GetLocation returns the physical location of the device
func (l Line) GetLocation() space.Vector {
	return l.Start
}

// SetLocation changes the physical location of the device
func (l *Line) SetLocation(newLocation space.Vector) {
	l.Start = newLocation
}

// GetOrientation returns the physical orientation of the device
func (l Line) GetOrientation() space.Orientation {
	return l.Direction
}

// SetOrientation changes the physical orientation of the device
func (l *Line) SetOrientation(newOrientation space.Orientation) {
	l.Direction = newOrientation
}
