package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Temporal is a Shifter which provides shifts that relate to changing time, Directionally
type Temporal struct {
	Start    *time.Time
	Interval *time.Duration
	Bender   ifaces.Bender
}

var _ ifaces.Shifter = (*Temporal)(nil)

// Shift returns a value representing some change or shift
func (s Temporal) Shift(t time.Time, l light.Light) float64 {
	secondsPast := float64(t.Sub(*s.Start)) / float64(*s.Interval)
	bend := s.Bender.Bend(secondsPast)
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Temporal) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.Start == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			t := p.Start()
			s.Start = &t
		})
	}
	if s.Interval == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Interval = p.SelectDuration()
		})
	}
	if s.Bender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Bender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.Bender.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Temporal) String() string {
	return fmt.Sprintf("shifter.Temporal{Start:%v, Interval:%v, Bender:%v}", s.Start, s.Interval, s.Bender)
}
