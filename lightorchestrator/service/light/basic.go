package light

import (
	"github.com/jmbarzee/color"
	"github.com/jmbarzee/space"
)

// Basic represents a NeoPixel Light in a line
type Basic struct {
	Color        color.Color
	Position     int
	NumPositions int
	Location     space.Cartesian
	Orientation  space.Spherical
}

// GetColor returns the color of the light
func (l Basic) GetColor() color.Color {
	return l.Color
}

// SetColor changes the color of the light
func (l *Basic) SetColor(newColor color.Color) {
	l.Color = newColor
}

// GetPosition returns the position of the Light (in a string)
func (l Basic) GetPosition() (int, int) {
	return l.Position, l.NumPositions
}

// GetLocation returns the point in space where the Light is
func (l Basic) GetLocation() space.Cartesian {
	return l.Location
}

// GetOrientation returns the direction the Light points
func (l Basic) GetOrientation() space.Spherical {
	return l.Orientation
}
