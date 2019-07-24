package lightdevice

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/placement"
)

type (
	NeoPixelBar struct {
		Start placement.Point
		Theta float32
		Phi   float32
	}
)

func NewNeoPixelBar() *NeoPixelBar {
	return &NeoPixelBar{}
}
