package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// Static is a Shifter which provides shifts that relate to changing time, Directionally
type Static struct {
	TheShift *float64
}

var _ ifaces.Shifter = (*Static)(nil)

// Shift returns a value representing some change or shift
func (s Static) Shift(t time.Time, l ifaces.Light) float64 {
	return *s.TheShift
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Static) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.TheShift == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.TheShift = p.SelectShift()
		})
	}
	return sFuncs
}

func (s Static) String() string {
	return fmt.Sprintf("shifter.Static{TheShift:%v}", s.TheShift)
}
