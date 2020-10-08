package painter

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service//effect/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/color"
)

// Static is a Painter which provides unchangeing colors
type Static struct {
	Color *color.HSLA
}

// GetColor returns a color based on t
func (p Static) GetColor(t time.Time) color.HSLA {
	return *p.Color
}

// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
func (p *Static) GetStabilizeFuncs() []shared.StabilizeFunc {
	sFuncs := []shared.StabilizeFunc{}
	if p.Color == nil {
		sFuncs = append(sFuncs, func(palette shared.Palette) error {
			c, err := palette.SelectColor()
			if err != nil {
				return err
			}
			p.Color = &c
			return nil
		})
	}
	return sFuncs
}
