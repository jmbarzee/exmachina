package neopixel

import (
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

const (
	neoPixelBarLength = 30
)

type Bar struct {
	device.BasicDevice
	*Line
}

func NewBar(uuid string, start space.Vector, direction space.Orientation) Bar {
	return Bar{
		BasicDevice: device.BasicDevice{
			ID: uuid,
		},
		Line: NewLine(start, direction, neoPixelBarLength),
	}
}

// GetType returns the type
func (d Bar) GetType() string {
	return "npBar"
}
