package ifaces

import (
	"time"

	"github.com/jmbarzee/color"
)

type Palette interface {
	Span

	// SelectBender returns a Bender
	SelectBender() Bender
	// SelectColor returns a Color
	SelectColor() color.Color
	// SelectDuration returns a Duration
	// Should generally range from 0.1s to 10s
	SelectDuration() *time.Duration
	// SelectShift returns a Shift
	// Should generally range from .01 to 1
	SelectShift() *float64
	// SelectShifter returns a Shifter
	SelectShifter() Shifter
	// SelectPainter returns a Painter
	SelectPainter() Painter
	// SelectEffect returns a Effect
	SelectEffect() Effect
}
