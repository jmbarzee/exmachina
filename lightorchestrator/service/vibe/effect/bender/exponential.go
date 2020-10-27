package bender

import (
	"fmt"
	"math"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Exponential is a Bender which provides a single unchanging bend
type Exponential struct {
	Exponent    *float64
	Coefficient *float64
}

var _ ifaces.Bender = (*Exponential)(nil)

// Bend returns a value representing some change or bend
func (s Exponential) Bend(f float64) float64 {
	bend := *s.Coefficient * math.Pow(math.Abs(f), *s.Exponent)
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Exponential) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.Exponent == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Exponent = p.SelectShift()
		})
	}
	if s.Coefficient == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Coefficient = p.SelectShift()
		})
	}
	return sFuncs
}

func (s Exponential) String() string {
	return fmt.Sprintf("shifter.Exponential{Exponent:%v, Coefficient:%v}", s.Exponent, s.Coefficient)
}
