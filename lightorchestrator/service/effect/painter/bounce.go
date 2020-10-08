package painter

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
	repeat "github.com/jmbarzee/services/lightorchestrator/service/shared"
)

// Bounce is a Painter which provides produces colors bouncing between ColorStart and ColorEnd,
// starting at p.ColorStart and shifting in the direction specified by Up
type Bounce struct {
	ColorStart *color.HSLA
	ColorEnd   *color.HSLA
	Up         *bool
	Shifter    shared.Shifter
}

// GetColor returns a color based on t
func (p Bounce) GetColor(t time.Time) color.HSLA {
	if *p.Up {
		if p.ColorStart.H < p.ColorEnd.H {
			hDistance := p.ColorEnd.H - p.ColorStart.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := p.ColorStart.H - p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(-hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	} else {
		if p.ColorStart.H > p.ColorEnd.H {
			hDistance := p.ColorStart.H - p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := (1 - p.ColorStart.H) + p.ColorEnd.H
			sDistance := p.ColorStart.S - p.ColorEnd.S
			lDistance := p.ColorStart.L - p.ColorEnd.L
			totalShift := p.Shifter.Shift(t)
			bounces := int(totalShift / hDistance)
			remainingShift := float32(math.Mod(float64(totalShift), float64(hDistance)))

			var hShift float32
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hDistance / hShift)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := *p.ColorStart
			c.ShiftHue(-hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	}
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Bounce) GetStabilizeFuncs() []shared.StabilizeFunc {
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
	if p.Up == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			b := repeat.RepeatableChance(palette.StartTime, .5)
			p.Up = &b
			return nil
		})
	}
	if p.ColorEnd == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			c, err := palette.SelectColor()
			if err != nil {
				return err
			}
			p.ColorEnd = &c
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
