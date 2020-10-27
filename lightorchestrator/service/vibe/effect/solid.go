package effect

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Solid is an Effect which displays all lights as a single color
type Solid struct {
	BasicEffect
	Painter ifaces.Painter
}

var _ ifaces.Effect = (*Solid)(nil)

// Render will produce a slice of lights based on the time and properties of lights
func (e Solid) Render(t time.Time, lights []light.Light) []light.Light {
	if len(lights) == 0 {
		return lights
	}

	c := e.Painter.Paint(t, lights[0])
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
func (e Solid) String() string {
	return fmt.Sprintf("effect.Solid{StartTime:%v, EndTime:%v, Rank:%v, Painter:%v}", e.StartTime, e.EndTime, e.Rank, e.Painter)
}
