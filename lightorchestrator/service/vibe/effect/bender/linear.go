package bender

import (
	"fmt"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Linear is a Bender which provides a single unchanging bend
type Linear struct {
	Rate *float64
}

var _ ifaces.Bender = (*Linear)(nil)

// Bend returns a value representing some change or bend
func (s Linear) Bend(f float64) float64 {
	bend := *s.Rate * f
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Linear) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.Rate == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Rate = p.SelectShift()
		})
	}
	return sFuncs
}

func (s Linear) String() string {
	return fmt.Sprintf("shifter.Linear{Rate:%v}", s.Rate)
}
