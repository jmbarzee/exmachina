package painter

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Rotate is a Painter which provides shifting colors starting at colorStart
type Rotate struct {
	ColorStart *color.HSLA
	Shifter    ifaces.Shifter
}

// Paint returns a color based on t
func (p Rotate) Paint(t time.Time) color.HSLA {
	newColor := *p.ColorStart
	newColor.ShiftHue(p.Shifter.Shift(t))
	return newColor
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Rotate) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if p.ColorStart == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.ColorStart = pa.SelectColor()
		})
	}
	if p.Shifter == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.Shifter = pa.SelectShifter()
		})
	} else {
		sFuncs = append(sFuncs, p.Shifter.GetStabilizeFuncs()...)
	}
	return sFuncs
}
