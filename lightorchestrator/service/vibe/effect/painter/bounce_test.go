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

func TestBouncePaint(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second
	theTruth := true
	theFalsehood := false
	aFloat := 1.0

	cases := []PainterTest{
		{
			Name: "Bounce from Green to Cyan ",
			Painter: &Bounce{
				ColorStart: &color.Green,
				ColorEnd:   &color.Cyan,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aSecond,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(time.Second * 0 / 24),
					ExpectedColor: color.Green,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 24),
					ExpectedColor: color.CoolGreen,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 24),
					ExpectedColor: color.GreenCyan,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 24),
					ExpectedColor: color.WarmCyan,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 24),
					ExpectedColor: color.Cyan,
				},
				{
					Time:          aTime.Add(time.Second * 5 / 24),
					ExpectedColor: color.WarmCyan,
				},
				{
					Time:          aTime.Add(time.Second * 6 / 24),
					ExpectedColor: color.GreenCyan,
				},
				{
					Time:          aTime.Add(time.Second * 7 / 24),
					ExpectedColor: color.CoolGreen,
				},
				{
					Time:          aTime.Add(time.Second * 8 / 24),
					ExpectedColor: color.Green,
				},
				{
					Time:          aTime.Add(time.Second * 9 / 24),
					ExpectedColor: color.CoolGreen,
				},
			},
		},
		{
			Name: "Bounce from Cyan to Green ",
			Painter: &Bounce{
				ColorStart: &color.Cyan,
				ColorEnd:   &color.Green,
				Up:         &theFalsehood,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aSecond,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(time.Second * 0 / 24),
					ExpectedColor: color.Cyan,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 24),
					ExpectedColor: color.WarmCyan,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 24),
					ExpectedColor: color.GreenCyan,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 24),
					ExpectedColor: color.CoolGreen,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 24),
					ExpectedColor: color.Green,
				},
				{
					Time:          aTime.Add(time.Second * 5 / 24),
					ExpectedColor: color.CoolGreen,
				},
				{
					Time:          aTime.Add(time.Second * 6 / 24),
					ExpectedColor: color.GreenCyan,
				},
				{
					Time:          aTime.Add(time.Second * 7 / 24),
					ExpectedColor: color.WarmCyan,
				},
				{
					Time:          aTime.Add(time.Second * 8 / 24),
					ExpectedColor: color.Cyan,
				},
				{
					Time:          aTime.Add(time.Second * 9 / 24),
					ExpectedColor: color.WarmCyan,
				},
			},
		},
		{
			Name: "Bounce from Orange to RedMagenta ",
			Painter: &Bounce{
				ColorStart: &color.Orange,
				ColorEnd:   &color.RedMagenta,
				Up:         &theFalsehood,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aSecond,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(time.Second * 0 / 24),
					ExpectedColor: color.Orange,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 24),
					ExpectedColor: color.WarmRed,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 24),
					ExpectedColor: color.Red,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 24),
					ExpectedColor: color.CoolRed,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 24),
					ExpectedColor: color.RedMagenta,
				},
				{
					Time:          aTime.Add(time.Second * 5 / 24),
					ExpectedColor: color.CoolRed,
				},
				{
					Time:          aTime.Add(time.Second * 6 / 24),
					ExpectedColor: color.Red,
				},
				{
					Time:          aTime.Add(time.Second * 7 / 24),
					ExpectedColor: color.WarmRed,
				},
				{
					Time:          aTime.Add(time.Second * 8 / 24),
					ExpectedColor: color.Orange,
				},
				{
					Time:          aTime.Add(time.Second * 9 / 24),
					ExpectedColor: color.WarmRed,
				},
			},
		},
		{
			Name: "Bounce from RedMagenta to Orange ",
			Painter: &Bounce{
				ColorStart: &color.RedMagenta,
				ColorEnd:   &color.Orange,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aSecond,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(time.Second * 0 / 24),
					ExpectedColor: color.RedMagenta,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 24),
					ExpectedColor: color.CoolRed,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 24),
					ExpectedColor: color.Red,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 24),
					ExpectedColor: color.WarmRed,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 24),
					ExpectedColor: color.Orange,
				},
				{
					Time:          aTime.Add(time.Second * 5 / 24),
					ExpectedColor: color.WarmRed,
				},
				{
					Time:          aTime.Add(time.Second * 6 / 24),
					ExpectedColor: color.Red,
				},
				{
					Time:          aTime.Add(time.Second * 7 / 24),
					ExpectedColor: color.CoolRed,
				},
				{
					Time:          aTime.Add(time.Second * 8 / 24),
					ExpectedColor: color.RedMagenta,
				},
				{
					Time:          aTime.Add(time.Second * 9 / 24),
					ExpectedColor: color.CoolRed,
				},
			},
		},
	}
	RunPainterTests(t, cases)
}

func TestBounceGetStabilizeFuncs(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aDuration := time.Second
	theTruth := true
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Bounce{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Bounce{
				ColorStart: &color.Red,
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
				Shifter:    &shifter.Temporal{},
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start: &aTime,
				},
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
				},
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
					Bender:   &bender.Linear{},
				},
			},
			&Bounce{
				ColorStart: &color.Red,
				ColorEnd:   &color.Red,
				Up:         &theTruth,
				Shifter: &shifter.Temporal{
					Start:    &aTime,
					Interval: &aDuration,
					Bender: &bender.Linear{
						Interval: &aFloat,
					},
				},
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Bender:   &bender.Linear{},
			Duration: aDuration,
			Color:    color.Red,
			Shift:    aFloat,
			Shifter:  &shifter.Temporal{},
		},
	}
	helper.RunStabilizeableTest(t, c)
}
