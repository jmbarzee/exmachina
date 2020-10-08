package shifter

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
)

// Sinusoidal is a Shifter which provides shifts that relate to changing time, linearly
type Sinusoidal struct {
	Start         *time.Time
	TimePerCycle  *time.Duration
	ShiftPerCycle *float32
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
func (s *Sinusoidal) GetStabilizeFuncs() []shared.StabilizeFunc {
	sFuncs := []shared.StabilizeFunc{}
	if s.Start == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			t := palette.StartTime
			s.Start = &t
			return nil
		})
	}
	if s.TimePerCycle == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			d, err := palette.SelectDuration()
			if err != nil {
				return err
			}
			s.TimePerCycle = &d
			return nil
		})
	}
	if s.ShiftPerCycle == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			shift, err := palette.SelectShift()
			if err != nil {
				return err
			}
			s.ShiftPerCycle = &shift
			return nil
		})
	}
	return sFuncs
}
