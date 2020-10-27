package painter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestMovePaint(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second

	cases := []PainterTest{
		{
			Name: "Paint Black",
			Painter: &Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Linear{
					Start:           &aTime,
					TimePerOneShift: &aSecond,
				},
			},
			Instants: func() []Instant {
				insts := make([]Instant, len(color.AllColors))
				for i := range insts {
					insts[i] = Instant{
						Time:          aTime.Add(time.Second * time.Duration(i) / 24),
						ExpectedColor: color.AllColors[i],
					}
				}
				return insts
			}(),
		},
	}
	RunPainterTests(t, cases)
}

func TestMoveGetStabilizeFuncs(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aDuration := time.Second
	c := helper.StabilizeableTest{
		Stabalizable: &Move{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Move{
				ColorStart: &color.Red,
			},
			&Move{
				ColorStart: &color.Red,
				Shifter:    &shifter.Linear{},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Linear{
					Start: &aTime,
				},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Linear{
					Start:           &aTime,
					TimePerOneShift: &aDuration,
				},
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Duration: aDuration,
			Color:    color.Red,
			Shifter:  &shifter.Linear{},
		},
	}
	helper.RunStabilizeableTest(t, c)
}
