package bender

import (
	"fmt"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Linear is a Bender which provides a single unchanging bend
type Linear struct {
	BendPerInput *float64
}

var _ ifaces.Bender = (*Linear)(nil)

// Bend returns a value representing some change or bend
func (s Linear) Bend(f float64) float64 {
	bend := f / *s.BendPerInput
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Linear) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.BendPerInput == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.BendPerInput = p.SelectShift()
		})
	}
	return sFuncs
}

func (s Linear) String() string {
	return fmt.Sprintf("shifter.Linear{BendPerInput:%v}", s.BendPerInput)
}
