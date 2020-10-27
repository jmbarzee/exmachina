package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Locational is a Shifter which provides shifts that relate to changing time, Directionally
type Locational struct {
	XBender ifaces.Bender
	YBender ifaces.Bender
	ZBender ifaces.Bender
}

var _ ifaces.Shifter = (*Locational)(nil)

// Shift returns a value representing some change or shift
func (s Locational) Shift(t time.Time, l light.Light) float64 {
	loc := l.GetLocation()
	bendX := s.XBender.Bend(loc.X)
	bendY := s.YBender.Bend(loc.Y)
	bendZ := s.ZBender.Bend(loc.Z)
	return bendX + bendY + bendZ
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Locational) GetStabilizeFuncs() []func(p ifaces.Palette) {
	sFuncs := []func(p ifaces.Palette){}
	if s.XBender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.XBender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.XBender.GetStabilizeFuncs()...)
	}
	if s.YBender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.YBender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.YBender.GetStabilizeFuncs()...)
	}
	if s.ZBender == nil {
		sFuncs = append(sFuncs, func(p ifaces.Palette) {
			s.ZBender = p.SelectBender()
		})
	} else {
		sFuncs = append(sFuncs, s.ZBender.GetStabilizeFuncs()...)
	}
	return sFuncs
}

func (s Locational) String() string {
	return fmt.Sprintf("shifter.Locational{XBender:%v, YBender:%v, ZBender:%v}", s.XBender, s.YBender, s.ZBender)
}
