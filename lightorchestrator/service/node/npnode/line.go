package neopixel

import (
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

// Line is a representation of a strait line of neopixels.
type Line struct {
	// Row provides the implementation of effect.Allocater
	*Row

	node.Basic

	space.Object
}

var _ node.Node = (*Line)(nil)

// NewLine creates a new Line
func NewLine(start space.Vector, direction, rotation space.Orientation, length int) *Line {

	d := &Line{
		Object: space.Object{
			Location:  start,
			Direction: direction,
			Rotation:  rotation,
		},
	}

	singleLEDVector := space.NewVector(direction, distPerLED)

	d.Row = NewRow(
		length,
		func() []light.Light {
			lights := make([]light.Light, length)
			for i := range lights {
				lightLocation := start.Translate(singleLEDVector.Scale(float64(i)))
				lightOrientation := rotation
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

// GetType returns the type
func (Line) GetType() string {
	return "NPLine"
}
