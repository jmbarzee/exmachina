package shifter

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/effect/shared"
)

const singleShift = 0.0000000001

// Linear is a Shifter which provides shifts that relate to changing time, linearly
type Linear struct {
	Start        *time.Time
	TimePerShift *time.Duration
}

// Shift returns a value representing some change or shift
func (s Linear) Shift(t time.Time) float32 {
	timePast := t.Sub(*s.Start)
	shifts := float32(timePast / *s.TimePerShift)
	return shifts * singleShift
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Linear) GetStabilizeFuncs() []shared.StabilizeFunc {
	sFuncs := []shared.StabilizeFunc{}
	if s.Start == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			t := palette.StartTime
			s.Start = &t
			return nil
		})
	}
	if s.TimePerShift == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			d, err := palette.SelectDuration()
			if err != nil {
				return err
			}
			s.TimePerShift = &d
			return nil
		})
	}
	return sFuncs
}
