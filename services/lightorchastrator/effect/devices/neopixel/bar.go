package neopixel

import "github.com/jmbarzee/domain/services/lightorchastrator/effect/space"

const (
	neoPixelBarLength = 60
)

func NewBar(
	start space.Vector,
	theta float64,
	phi float64,
) *Line {
	return NewLine(start, theta, phi, neoPixelBarLength)
}
