package painter

import (
	"fmt"
	"math"
	"time"

	"github.com/jmbarzee/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/repeatable"
)

// Bounce is a Painter which provides produces colors bouncing between ColorStart and ColorEnd,
// starting at p.ColorStart and shifting in the direction specified by Up
type Bounce struct {
	ColorStart color.Color
	ColorEnd   color.Color
	Up         *bool
	Shifter    ifaces.Shifter
}

var _ ifaces.Painter = (*Bounce)(nil)

// Paint returns a color based on t
func (p Bounce) Paint(t time.Time, l ifaces.Light) color.Color {
	start := p.ColorStart.HSL()
	end := p.ColorEnd.HSL()
	if *p.Up {
		if start.H < end.H {
			hDistance := end.H - start.H
			sDistance := start.S - end.S
			lDistance := start.L - end.L
			totalShift := p.Shifter.Shift(t, l)
			bounces := int(totalShift / hDistance)
			remainingShift := math.Mod(totalShift, hDistance)

			var hShift float64
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = remainingShift
			} else {
				// odd number of bounces
				hShift = hDistance - remainingShift
			}
			hShiftRatio := (hShift / hDistance)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := start
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := (1 - start.H) + end.H
			sDistance := start.S - end.S
			lDistance := start.L - end.L
			totalShift := p.Shifter.Shift(t, l)
			bounces := int(totalShift / hDistance)
			remainingShift := math.Mod(totalShift, hDistance)

			var hShift float64
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = -remainingShift
			} else {
				// odd number of bounces
				hShift = -(hDistance - remainingShift)
			}
			hShiftRatio := (hShift / hDistance)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := start
			c.ShiftHue(-hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	} else {
		if start.H > end.H {
			hDistance := start.H - end.H
			sDistance := start.S - end.S
			lDistance := start.L - end.L
			totalShift := p.Shifter.Shift(t, l)
			bounces := int(totalShift / hDistance)
			remainingShift := math.Mod(totalShift, hDistance)

			var hShift float64
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = -remainingShift
			} else {
				// odd number of bounces
				hShift = -(hDistance - remainingShift)
			}
			hShiftRatio := (hShift / hDistance)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := start
			c.ShiftHue(hShift)
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		} else {
			hDistance := start.H + (1 - end.H)
			sDistance := start.S - end.S
			lDistance := start.L - end.L
			totalShift := p.Shifter.Shift(t, l)
			bounces := int(totalShift / hDistance)
			remainingShift := math.Mod(totalShift, hDistance)

			var hShift float64
			if (bounces % 2) == 0 {
				// even number of bounces
				hShift = -remainingShift
			} else {
				// odd number of bounces
				hShift = -(hDistance - remainingShift)
			}
			hShiftRatio := (hShift / hDistance)
			sShift := sDistance * hShiftRatio
			lShift := lDistance * hShiftRatio

			c := start
			c.ShiftHue(hShift) // shifting past 0
			c.SetSaturation(c.S + sShift)
			c.SetLightness(c.L + lShift)

			return c
		}
	}
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Bounce) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if p.ColorStart == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.ColorStart = pa.SelectColor().HSL()
		})
	}
	if p.ColorEnd == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			p.ColorEnd = pa.SelectColor().HSL()
		})
	}
	if p.Up == nil {
		sFuncs = append(sFuncs, func(pa ifaces.Palette) {
			b := repeatable.Chance(pa.Start(), .5)
			p.Up = &b
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

func (p Bounce) String() string {
	return fmt.Sprintf("painter.Bounce{ColorStart:%v, ColorEnd:%v, Up:%v, Shifter:%v}", p.ColorStart, p.ColorEnd, p.Up, p.Shifter)
}
