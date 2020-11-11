package bender

import (
	"fmt"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// Static is a Bender which provides a single unchanging bend
type Static struct {
	TheBend *float64
}

var _ ifaces.Bender = (*Static)(nil)

// Bend returns a value representing some change or bend
func (s Static) Bend(f float64) float64 {
	return *s.TheBend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Static) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.TheBend == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.TheBend = p.SelectShift()
		})
	}
	return sFuncs
}

func (s Static) String() string {
	return fmt.Sprintf("shifter.Static{TheBend:%v}", s.TheBend)
}
