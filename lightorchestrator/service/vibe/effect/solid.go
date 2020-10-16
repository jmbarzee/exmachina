package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Solid is an Effect which displays all lights as a single color
type Solid struct {
	BasicEffect
	Painter ifaces.Painter
}

// Render will produce a slice of lights based on the time and properties of lights
func (e Solid) Render(t time.Time, lights []light.Light) []light.Light {
	c := e.Painter.Paint(t)
	for i := range lights {
		lights[i].SetColor(c)
	}
	return lights
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (e *Solid) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if e.Painter == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			e.Painter = pa.SelectPainter()
		})
	} else {
		sFuncs = append(sFuncs, e.Painter.GetStabilizeFuncs()...)
	}
	return sFuncs
}
