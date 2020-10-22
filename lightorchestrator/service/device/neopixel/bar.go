package neopixel

import (
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

const (
	npBarLength = 2

	ledsPerNPBar = npBarLength * ledsPerMeter
)

type Bar struct {
	device.BasicDevice
	*Line
}

func NewBar(uuid string, start space.Vector, direction, rotation space.Orientation) Bar {
	return Bar{
		BasicDevice: device.BasicDevice{
			ID: uuid,
		},
		Line: NewLine(start, direction, rotation, ledsPerNPBar),
	}
}

// GetType returns the type
func (d Bar) GetType() string {
	return "npBar"
}
