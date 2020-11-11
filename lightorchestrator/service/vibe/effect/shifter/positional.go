package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// Positional is a Shifter which provides shifts that relate to changing time, Directionally
type Positional struct {
	Bender ifaces.Bender
}

var _ ifaces.Shifter = (*Positional)(nil)

// Shift returns a value representing some change or shift
func (s Positional) Shift(t time.Time, l ifaces.Light) float64 {
	pos, numPos := l.GetPosition()
	bend := s.Bender.Bend(float64(pos) / float64(numPos))
	return bend
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Positional) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.Bender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.Bender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.Bender.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Positional) String() string {
	return fmt.Sprintf("shifter.Positional{Bender:%v}", s.Bender)
}
