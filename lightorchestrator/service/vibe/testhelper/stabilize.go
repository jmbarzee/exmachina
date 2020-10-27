package testhelper

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
)

type StabilizeableTest struct {
	Stabalizable     ifaces.Stabalizable
	ExpectedVersions []ifaces.Stabalizable
	Palette          TestPalette
}

func RunStabilizeableTest(t *testing.T, c StabilizeableTest) {
	for i, expectedVersion := range c.ExpectedVersions {

		sFuncs := c.Stabalizable.GetStabilizeFuncs()
		if len(sFuncs) > len(c.ExpectedVersions)-i {
			t.Fatalf("Stabalizable %v failed. Unexpected number of sFuncs:\n\tExpected: %v,\n\tActual: %v", i, len(c.ExpectedVersions)-i, len(sFuncs))
		}
		sFuncs[0](c.Palette)
		if !StructsEqual(c.Stabalizable, expectedVersion) {
			t.Fatalf("Stabalizable %v failed. Stabilizables were not equal:\n\tExpected: %v,\n\tActual: %v", i, expectedVersion, c.Stabalizable)
		}
	}
}

type TestPalette struct {
	span.Span
	Bender   ifaces.Bender
	Color    color.HSLA
	Duration time.Duration
	Shift    float64
	Shifter  ifaces.Shifter
	Painter  ifaces.Painter
	Effect   ifaces.Effect
}

// SelectBender returns a Bender
func (p TestPalette) SelectBender() ifaces.Bender {
	return p.Bender
}

// SelectColor returns a Color
func (p TestPalette) SelectColor() *color.HSLA {
	return &p.Color
}

// SelectDuration returns a Duration
// Should generally range from 0.1s to 10s
func (p TestPalette) SelectDuration() *time.Duration {
	return &p.Duration
}

// SelectShift returns a Shift
// Should generally range from .01 to 1
func (p TestPalette) SelectShift() *float64 {
	return &p.Shift
}

// SelectShifter returns a Shifter
func (p TestPalette) SelectShifter() ifaces.Shifter {
	return p.Shifter
}

// SelectPainter returns a Painter
func (p TestPalette) SelectPainter() ifaces.Painter {
	return p.Painter
}

// SelectEffect returns a Effect
func (p TestPalette) SelectEffect() ifaces.Effect {
	return p.Effect
}
