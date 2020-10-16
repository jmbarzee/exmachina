package vibe

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/repeatable"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/painter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
)

// Basic is a vibe which can produce most Effects
type Basic struct {
	span.Span
	Effects []ifaces.Effect
}

// Stabilize locks in part of the visual representation of a vibe.
func (v *Basic) Stabilize() ifaces.Vibe {
	sFuncs := v.GetStabilizeFuncs()
	option := repeatable.Option(v.Start(), len(sFuncs))
	sFuncs[option](v)
	return v
}

// Materialize locks all remaining unlocked visuals of a vibe
// then returns the resulting effects
func (v *Basic) Materialize() []ifaces.Effect {
	for {
		sFuncs := v.GetStabilizeFuncs()
		for _, sFunc := range sFuncs {
			sFunc(v)
		}
	}
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (v *Basic) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	for _, e := range v.Effects {
		sFuncs = append(sFuncs, e.GetStabilizeFuncs()...)
	}
	sFuncs = append(sFuncs, func(p ifaces.Palette) {
		v.Effects = append(v.Effects, v.SelectEffect())
	})
	return sFuncs
}

// ifaces.Palette implementation

// SelectColor returns a Color
func (v Basic) SelectColor() *color.HSLA {
	length := len(color.AllColors)
	option := repeatable.Option(v.Start(), length)
	c := color.AllColors[option]
	return &c
}

// SelectDuration returns a Duration
func (v Basic) SelectDuration() *time.Duration {
	min := time.Second / 10
	max := time.Second * 10
	d := repeatable.RandDuration(v.Start(), min, max)
	return &d
}

// SelectShift returns a Shift
func (v Basic) SelectShift() *float32 {
	min := float32(0.01)
	max := float32(1.00)
	s := repeatable.RandShift(v.Start(), min, max, 0.001)
	return &s

}

// SelectShifter returns a Shifter
func (v Basic) SelectShifter() ifaces.Shifter {
	options := []ifaces.Shifter{
		&shifter.Linear{},
		&shifter.Sinusoidal{},
	}
	length := len(options)
	option := repeatable.Option(v.Start(), length)

	return options[option]
}

// SelectPainter returns a Painter
func (v Basic) SelectPainter() ifaces.Painter {
	options := []ifaces.Painter{
		&painter.Static{},
		&painter.Rotate{},
		&painter.Bounce{},
	}
	length := len(options)
	option := repeatable.Option(v.Start(), length)

	return options[option]
}

// SelectEffect returns a Effect
func (v Basic) SelectEffect() ifaces.Effect {
	options := []ifaces.Effect{
		&effect.Solid{},
		&effect.Future{},
	}
	length := len(options)
	option := repeatable.Option(v.Start(), length)

	return options[option]
}
