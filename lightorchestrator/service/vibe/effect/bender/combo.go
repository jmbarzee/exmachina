package bender

import (
	"fmt"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Combo is a Bender which provides a single unchanging bend
type Combo struct {
	A ifaces.Bender
	B ifaces.Bender
}

var _ ifaces.Bender = (*Combo)(nil)

// Bend returns a value representing some change or bend
func (s Combo) Bend(f float64) float64 {
	bend := s.A.Bend(f) + s.B.Bend(f)
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Combo) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.A == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.A = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.A.GetStabilizeFuncs()...)
	}
	if s.B == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.B = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.B.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Combo) String() string {
	return fmt.Sprintf("shifter.Combo{A:%v, B:%v}", s.A, s.B)
}
