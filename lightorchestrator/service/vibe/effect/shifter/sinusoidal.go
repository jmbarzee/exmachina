package shifter

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Sinusoidal is a Shifter which provides shifts that relate to changing time, linearly
type Sinusoidal struct {
	Start         *time.Time
	TimePerCycle  *time.Duration // Period
	ShiftPerCycle *float32       // Amplitude
}

// Shift returns a value representing some change or shift
func (s Sinusoidal) Shift(t time.Time) float32 {
	timePast := t.Sub(*s.Start)
	cycles := timePast / *s.TimePerCycle
	sin := float32(math.Sin(2 * math.Pi * float64(cycles)))
	normalizedSin := (sin + 1) / 2 // normalized sin has the correct frequency but y ranges from 0 to 1
	return *s.ShiftPerCycle * normalizedSin
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Sinusoidal) GetStabilizeFuncs() []func(ifaces.Palette) {
	sFuncs := []func(ifaces.Palette){}
	if s.Start == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			t := p.Start()
			s.Start = &t
		})
	}
	if s.TimePerCycle == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.TimePerCycle = p.SelectDuration()
		})
	}
	if s.ShiftPerCycle == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.ShiftPerCycle = p.SelectShift()
		})
	}
	return sFuncs
}
