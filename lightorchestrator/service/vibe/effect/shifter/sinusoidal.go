package shifter

import (
	"fmt"
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Sinusoidal is a Shifter which provides shifts that relate to changing time, linearly
type Sinusoidal struct {
	Start         *time.Time
	TimePerCycle  *time.Duration // Period
	ShiftPerCycle *float64       // Amplitude
}

var _ ifaces.Shifter = (*Sinusoidal)(nil)

// Shift returns a value representing some change or shift
func (s Sinusoidal) Shift(t time.Time, l light.Light) float64 {
	timePast := t.Sub(*s.Start)
	cycles := float64(timePast) / float64(*s.TimePerCycle)
	sin := math.Sin(2 * math.Pi * cycles)
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

func (s Sinusoidal) String() string {
	return fmt.Sprintf("shifter.Sinusoidal{Start:%v, TimePerCycle:%v, ShiftPerCycle:%v}", s.Start, s.TimePerCycle, s.ShiftPerCycle)
}
