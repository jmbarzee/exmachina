package effects

import (
	"time"

	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/color"
)

type (
	Solid struct {
		BasicEffect
		Color ColorFunc
	}
)

func (e Solid) Render(t time.Time, lights []effect.Light) []color.HSLA {
	c := e.Color(t)
	colors := make([]color.HSLA, len(lights))
	for i := range colors {
		colors[i] = c
	}
	return colors
}
