package ifaces

import (
	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/space"
)

type Light interface {

	// GetColor returns the color of the light
	GetColor() color.HSLA
	// SetColor changes the color of the light
	SetColor(newColor color.HSLA)

	// GetPosition returns the position of the Light (in a string)
	GetPosition() (int, int)
	// GetLocation returns the point in space where the Light is
	GetLocation() space.Cartesian
	// GetOrientation returns the direction the Light points
	GetOrientation() space.Spherical
}
