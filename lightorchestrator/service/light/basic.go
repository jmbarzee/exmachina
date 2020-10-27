package light

import (
	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

// Basic represents a NeoPixel Light in a line
type Basic struct {
	Color        color.HSLA
	Position     int
	NumPositions int
	Location     space.Vector
	Orientation  space.Orientation
}

// GetColor returns the color of the light
func (l Basic) GetColor() color.HSLA {
	return l.Color
}

// SetColor changes the color of the light
func (l *Basic) SetColor(newColor color.HSLA) {
	l.Color = newColor
}

// GetPosition returns the position of the Light (in a string)
func (l Basic) GetPosition() (int, int) {
	return l.Position, l.NumPositions
}

// GetLocation returns the point in space where the Light is
func (l Basic) GetLocation() space.Vector {
	return l.Location
}

// GetOrientation returns the direction the Light points
func (l Basic) GetOrientation() space.Orientation {
	return l.Orientation
}
