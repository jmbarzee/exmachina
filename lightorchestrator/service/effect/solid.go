package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/shared"
)

type (
	Solid struct {
		BasicEffect
		Color ColorFunc
	}
)

func (e Solid) Render(t time.Time, lights []shared.Light) []shared.Light {
	c := e.Color(t)
	for i := range lights {
		lights[i].Color = c
	}
	return lights
}
