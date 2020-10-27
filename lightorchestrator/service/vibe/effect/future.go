package effect

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Future is an Effect which displays each consecutive light
// as the "future" of the previous light
type Future struct {
	BasicEffect
	TimePerLight *time.Duration
	Painter      ifaces.Painter
}

var _ ifaces.Effect = (*Future)(nil)

// Render will produce a slice of lights based on the time and properties of lights
func (e Future) Render(t time.Time, lights []light.Light) []light.Light {
	for i, l := range lights {
		distanceInFuture := *e.TimePerLight * time.Duration(i)
		c := e.Painter.Paint(t.Add(distanceInFuture), l)
		lights[i].SetColor(c)
	}
	return lights
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (e *Future) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if e.TimePerLight == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			e.TimePerLight = pa.SelectDuration()
		})
	}
	if e.Painter == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			e.Painter = pa.SelectPainter()
		})
	} else {
		sFuncs = append(sFuncs, e.Painter.GetStabilizeFuncs()...)
	}

	return sFuncs
}

func (e Future) String() string {
	return fmt.Sprintf("effect.Future{StartTime:%v, EndTime:%v, Rank:%v, TimePerLight:%v, Painter:%v}", e.StartTime, e.EndTime, e.Rank, e.TimePerLight, e.Painter)
}
