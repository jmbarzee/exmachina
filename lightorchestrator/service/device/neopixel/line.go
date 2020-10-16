package neopixel

import (
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
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
	// Rotation is direction which all LEDs point (orthogonal to Direction)
	Rotation space.Orientation
}

// NewLine creates a new Line
func NewLine(start space.Vector, direction, rotation space.Orientation, length int) *Line {
	// TODO

	d := &Line{
		Start:     start,
		Direction: direction,
		Rotation:  rotation,
	}

	singleLEDVector := space.NewVector(direction, distPerLED)

	d.Row = NewRow(
		length,
		func() []light.Light {
			lights := make([]light.Light, length)
			for i := range lights {
				lightLocation := start.Translate(singleLEDVector.Scale(float32(i)))
				lightOrientation := rotation
				lights[i] = &light.Basic{
					Position:    i,
					Location:    lightLocation,
					Orientation: lightOrientation,
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
