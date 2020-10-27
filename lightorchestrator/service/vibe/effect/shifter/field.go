package shifter

import (
	"fmt"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// Field is a Shifter which provides shifts that relate to changing time, linearly
type Field struct {
	XBender ifaces.Bender
	YBender ifaces.Bender
	ZBender ifaces.Bender
}

var _ ifaces.Shifter = (*Field)(nil)

// Shift returns a value representing some change or shift
func (s Field) Shift(t time.Time, l light.Light) float64 {
	loc := l.GetLocation()
	shiftX := s.XBender.Bend(loc.X)
	shiftY := s.YBender.Bend(loc.Y)
	shiftZ := s.ZBender.Bend(loc.Z)
	return shiftX + shiftY + shiftZ
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (s *Field) GetStabilizeFuncs() []func(p ifaces.Palette) {
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

func (s Field) String() string {
	return fmt.Sprintf("shifter.Field{XBender:%v, TimePerOneShift:%v, TimePerOneShift:%v}", s.XBender, s.YBender, s.ZBender)
}
