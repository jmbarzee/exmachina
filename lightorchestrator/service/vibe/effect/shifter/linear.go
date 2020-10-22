package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// OneShift is just 1
// it can represent a full wrap around of Hue or something else
const OneShift = 1.0

// Linear is a Shifter which provides shifts that relate to changing time, linearly
type Linear struct {
	Start           *time.Time
	TimePerOneShift *time.Duration
}

var _ ifaces.Shifter = (*Linear)(nil)

// Shift returns a value representing some change or shift
func (s Linear) Shift(t time.Time) float64 {
	timePast := t.Sub(*s.Start)
	shift := float64(timePast) / float64(*s.TimePerOneShift)
	return shift * OneShift
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Linear) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.Start == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			t := p.Start()
			s.Start = &t
		})
	}
	if s.TimePerOneShift == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.TimePerOneShift = p.SelectDuration()
		})
	}
	return sFuncs
}

func (s Linear) String() string {
	return fmt.Sprintf("shifter.Linear{Start:%v, TimePerOneShift:%v}", s.Start, s.TimePerOneShift)
}
