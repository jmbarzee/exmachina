package painter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

type (
	PaintTest struct {
		Name     string
		Painter  ifaces.Painter
		Instants []Instant
	}

	Instant struct {
		Time          time.Time
		ExpectedColor color.HSLA
	}
)

func RunPainterTests(t *testing.T, cases []PaintTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			for i, instant := range c.Instants {
				actualColor := c.Painter.Paint(instant.Time)
				if !helper.ColorsEqual(instant.ExpectedColor, actualColor) {
					t.Fatalf("instant %v failed:\n\tExpected: %v,\n\tActual: %v", i, instant.ExpectedColor, actualColor)
				}
			}
		})
	}
}
