package painter

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Static is a Painter which provides unchangeing colors
type Static struct {
	Color *color.HSLA
}

// Paint returns a color based on t
func (p Static) Paint(t time.Time) color.HSLA {
	return *p.Color
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Static) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if p.Color == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.Color = pa.SelectColor()
		})
	}
	return sFuncs
}
