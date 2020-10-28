package painter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestMovePaint(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second
	aFloat := 1.0

	cases := []PainterTest{
		{
			Name: "Paint all the colors",
			Painter: &Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aSecond,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
			Instants: func() []Instant {
				insts := make([]Instant, len(color.AllColors))
				for i := range insts {
					insts[i] = Instant{
						Time:          aTime.Add(aSecond * time.Duration(i) / 24),
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
	aFloat := 1.1
	aDuration := time.Second
	c := helper.StabilizeableTest{
		Stabalizable: &Move{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Move{
				ColorStart: &color.Red,
			},
			&Move{
				ColorStart: &color.Red,
				Shifter:    &shifter.Temporal{},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Temporal{
					Start: &aTime,
				},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
				},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
					Bender:   &bender.Static{},
				},
			},
			&Move{
				ColorStart: &color.Red,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
					Bender: &bender.Static{
						TheBend: &aFloat,
					},
				},
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Bender:   &bender.Static{},
			Duration: aDuration,
			Color:    color.Red,
			Shift:    aFloat,
			Shifter:  &shifter.Temporal{},
		},
	}
	helper.RunStabilizeableTest(t, c)
}
