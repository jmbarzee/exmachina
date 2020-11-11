package painter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// Move is a Painter which provides shifting colors starting at colorStart
type Move struct {
	ColorStart *color.HSLA
	Shifter    ifaces.Shifter
}

var _ ifaces.Painter = (*Move)(nil)

// Paint returns a color based on t
func (p Move) Paint(t time.Time, l ifaces.Light) color.HSLA {
	newColor := *p.ColorStart
	shift := p.Shifter.Shift(t, l)
	newColor.ShiftHue(shift)
	return newColor
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Move) GetStabilizeFuncs() []func(p ifaces.Palette) {
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

func (p Move) String() string {
	return fmt.Sprintf("painter.Move{ColorStart:%v, Shifter:%v}", p.ColorStart, p.Shifter)
}
