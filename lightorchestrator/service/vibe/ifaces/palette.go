package ifaces

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
)

type Palette interface {
	Span

	// SelectColor returns a Color
	SelectColor() *color.HSLA
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
