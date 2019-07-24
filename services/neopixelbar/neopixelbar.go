package neopixelbar

import (
	lightsub "github.com/jmbarzee/domain/services/lightorchastrator/client"
)

type (
	NeoBar struct {
		*lightsub.LightSub
	}
)

const (
	size = 30
)

func NewNeoBar(port int, domainPort int) NeoBar {
	return NeoBar{
		LightSub: lightsub.NewLightSub(port, size, "neoPixelBar", domainPort),
	}
}
