package painter

import (
	"testing"
	"time"

	"github.com/jmbarzee/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

type (
	PainterTest struct {
		Name     string
		Painter  ifaces.Painter
		Instants []Instant
	}

	Instant struct {
		Time          time.Time
		Light         ifaces.Light
		ExpectedColor color.Color
	}
)

func RunPainterTests(t *testing.T, cases []PainterTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			for i, instant := range c.Instants {
				actualColor := c.Painter.Paint(instant.Time, instant.Light)
				if !helper.ColorsEqual(instant.ExpectedColor, actualColor) {
					t.Fatalf("instant %v failed:\n\tExpected: %v,\n\tActual: %v", i, instant.ExpectedColor, actualColor)
				}
			}
		})
	}
}
