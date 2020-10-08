package shared

import (
	"errors"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/shared"
)

var (
	NeedLargerColorPaletteError    = errors.New("Require a larger Color palette to Stabilize")
	NeedLargerDurationPaletteError = errors.New("Require a larger Duration palette to Stabilize")
	NeedLargerShiftPaletteError    = errors.New("Require a larger Shift palette to Stabilize")
	NeedLargerShifterPaletteError  = errors.New("Require a larger Shifter palette to Stabilize")
	NeedLargerPainterPaletteError  = errors.New("Require a larger Painter palette to Stabilize")
)

type Palette struct {
	StartTime time.Time
	Colors    []color.HSLA
	Durations []time.Duration
	Shifts    []float32
	Shifters  []Shifter
	Painters  []Painter
}

// SelectColor returns a Color which has been removed from the Palette
func (p *Palette) SelectColor() (color.HSLA, error) {
	length := len(p.Colors)
	if length < 1 {
		return color.HSLA{}, NeedLargerColorPaletteError
	}
	option := shared.RepeatableOption(p.StartTime, length)
	c := p.Colors[option]
	p.Colors[option] = p.Colors[length-1]
	p.Colors = p.Colors[:length-1] // Truncate slice.
	return c, nil
}

// SelectDuration returns a Duration which has been removed from the Palette
func (p *Palette) SelectDuration() (time.Duration, error) {
	length := len(p.Durations)
	if length < 1 {
		return 0, NeedLargerDurationPaletteError
	}
	option := shared.RepeatableOption(p.StartTime, length)
	d := p.Durations[option]
	p.Durations[option] = p.Durations[length-1]
	p.Durations = p.Durations[:length-1] // Truncate slice.
	return d, nil
}

// SelectShift returns a Shift which has been removed from the Palette
func (p *Palette) SelectShift() (float32, error) {
	length := len(p.Shifts)
	if length < 1 {
		return 0, NeedLargerShiftPaletteError
	}
	option := shared.RepeatableOption(p.StartTime, length)
	s := p.Shifts[option]
	p.Shifts[option] = p.Shifts[length-1]
	p.Shifts = p.Shifts[:length-1] // Truncate slice.
	return s, nil
}

// SelectShifter returns a Shifter which has been removed from the Palette
func (p *Palette) SelectShifter() (Shifter, error) {
	length := len(p.Shifters)
	if length < 1 {
		return nil, NeedLargerShifterPaletteError
	}
	option := shared.RepeatableOption(p.StartTime, length)
	s := p.Shifters[option]
	p.Shifters[option] = p.Shifters[length-1]
	p.Shifters = p.Shifters[:length-1] // Truncate slice.
	return s, nil
}

// SelectPainter returns a Painter which has been removed from the Palette
func (p *Palette) SelectPainter() (Painter, error) {
	length := len(p.Painters)
	if length < 1 {
		return nil, NeedLargerPainterPaletteError
	}
	option := shared.RepeatableOption(p.StartTime, length)
	pa := p.Painters[option]
	p.Painters[option] = p.Painters[length-1]
	p.Painters = p.Painters[:length-1] // Truncate slice.
	return pa, nil
}
