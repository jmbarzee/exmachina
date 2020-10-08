package painter

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
)

// Rotate is a Painter which provides shifting colors starting at colorStart
type Rotate struct {
	ColorStart *color.HSLA
	Shifter    shared.Shifter
}

// GetColor returns a color based on t
func (p Rotate) GetColor(t time.Time) color.HSLA {
	newColor := *p.ColorStart
	newColor.ShiftHue(p.Shifter.Shift(t))
	return newColor
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Rotate) GetStabilizeFuncs() []shared.StabilizeFunc {
	sFuncs := []shared.StabilizeFunc{}
	if p.ColorStart == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			c, err := palette.SelectColor()
			if err != nil {
				return err
			}
			p.ColorStart = &c
			return nil
		})
	}
	if p.Shifter == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			shifter, err := palette.SelectShifter()
			if err != nil {
				return err
			}
			p.Shifter = shifter
			return nil
		})
	} else {
		sFuncs = append(sFuncs, p.Shifter.GetStabilizeFuncs()...)
	}
	return sFuncs
}
