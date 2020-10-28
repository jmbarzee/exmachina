package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Combo is a Shifter which provides shifts that relate to changing time, Directionally
type Combo struct {
	A ifaces.Shifter
	B ifaces.Shifter
}

var _ ifaces.Shifter = (*Combo)(nil)

// Shift returns a value representing some change or shift
func (s Combo) Shift(t time.Time, l light.Light) float64 {
	shift := s.A.Shift(t, l) + s.B.Shift(t, l)
	return shift
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Combo) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.A == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.A = p.SelectShifter()
		})
	} else {
		sFuncs = append(sFuncs, s.A.GetStabilizeFuncs()...)
	}
	if s.B == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.B = p.SelectShifter()
		})
	} else {
		sFuncs = append(sFuncs, s.B.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Combo) String() string {
	return fmt.Sprintf("shifter.Combo{A:%v, B:%v}", s.A, s.B)
}
