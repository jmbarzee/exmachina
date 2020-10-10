package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
)

// Future is an Effect which displays each consecutive light
// as the "future" of the previous light
type Future struct {
	BasicEffect
	Painter      shared.Painter
	TimePerLight time.Duration
}

// Render will produce a slice of lights based on the time and properties of lights
func (e Future) Render(t time.Time, lights []light.Light) []light.Light {
	for i := range lights {
		distanceInFuture := e.TimePerLight * time.Duration(i)
		c := e.Painter.GetColor(t.Add(distanceInFuture))
		lights[i].Color = c
	}
	return lights
}
