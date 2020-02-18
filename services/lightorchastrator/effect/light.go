package effect

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/effect/color"
	"github.com/jmbarzee/domain/services/lightorchastrator/effect/space"
)

// Light represents a NeoPixel in a line
type Light struct {
	Position        int
	Color           color.HSLA
	GetLocationFunc func(position int) space.Vector
}

// GetColor returns the color of the light
func (l Light) GetColor() color.HSLA {
	return l.Color
}

// SetColor changes the color of the light
func (l *Light) SetColor(newColor color.HSLA) {
	l.Color = newColor
}

// GetLocation returns the point in space where the Light is
func (l Light) GetLocation() space.Vector {
	return l.GetLocationFunc(l.Position)
}

// GetPosition returns the position of the Light (in a string)
func (l Light) GetPosition() int {
	return l.Position
}
