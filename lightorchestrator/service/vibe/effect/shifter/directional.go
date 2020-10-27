package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// OneShift is just 1
// it can represent a full wrap around of Hue or something else
const OneShift = 1.0

// Directional is a Shifter which provides shifts that relate to changing time, Directionally
type Directional struct {
	PhiBender   ifaces.Bender
	ThetaBender ifaces.Bender
}

var _ ifaces.Shifter = (*Directional)(nil)

// Shift returns a value representing some change or shift
func (s Directional) Shift(t time.Time, l light.Light) float64 {
	ori := l.GetOrientation()
	bendPhi := s.PhiBender.Bend(ori.Phi)
	bendTheta := s.ThetaBender.Bend(ori.Theta)
	return bendPhi + bendTheta
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Directional) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.PhiBender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.PhiBender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.PhiBender.GetStabilizeFuncs()...)
	}
	if s.ThetaBender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.ThetaBender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.ThetaBender.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Directional) String() string {
	return fmt.Sprintf("shifter.Directional{PhiBender:%v, ThetaBender:%v}", s.PhiBender, s.ThetaBender)
}
