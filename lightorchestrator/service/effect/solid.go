package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
)

// Solid is an Effect which displays all lights as a single color
type Solid struct {
	BasicEffect
	Painter shared.Painter
}

// Render will produce a slice of lights based on the time and properties of lights
func (e Solid) Render(t time.Time, lights []light.Light) []light.Light {
	c := e.Painter.GetColor(t)
	for i := range lights {
		lights[i].Color = c
	}
	return lights
}
