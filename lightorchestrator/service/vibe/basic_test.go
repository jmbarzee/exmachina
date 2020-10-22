package vibe

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/painter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestBasicStabilize(t *testing.T) {
	aTime1 := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aTime2 := aTime1.Add(time.Second)
	aTime3 := aTime2.Add(time.Minute)
	aDuration := time.Nanosecond * 2245197264
	aDuration2 := time.Nanosecond * 4291714913

	cases := []StabilizeTest{
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: span.Span{
					StartTime: aTime1,
				},
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: span.Span{
						StartTime: aTime1,
					},
					Effects: []ifaces.Effect{
						&effect.Future{},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime1,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							TimePerLight: &aDuration,
						},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime1,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							TimePerLight: &aDuration,
							Painter:      &painter.Move{},
						},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime1,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							TimePerLight: &aDuration,
							Painter: &painter.Move{
								ColorStart: &color.RedMagenta,
							},
						},
					},
				},
			},
		},
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: span.Span{
					StartTime: aTime2,
				},
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: span.Span{
						StartTime: aTime2,
					},
					Effects: []ifaces.Effect{
						&effect.Solid{},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime2,
					},
					Effects: []ifaces.Effect{
						&effect.Solid{
							Painter: &painter.Static{},
						},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime2,
					},
					Effects: []ifaces.Effect{
						&effect.Solid{
							Painter: &painter.Static{
								Color: &color.Green,
							},
						},
					},
				},
			},
		},
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: span.Span{
					StartTime: aTime3,
				},
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: span.Span{
						StartTime: aTime3,
					},
					Effects: []ifaces.Effect{
						&effect.Future{},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime3,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							Painter: &painter.Static{},
						},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime3,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							Painter: &painter.Static{
								Color: &color.WarmGreen,
							},
						},
					},
				},
				&Basic{
					Span: span.Span{
						StartTime: aTime3,
					},
					Effects: []ifaces.Effect{
						&effect.Future{
							TimePerLight: &aDuration2,
							Painter: &painter.Static{
								Color: &color.WarmGreen,
							},
						},
					},
				},
			},
		},
	}

	RunStabilizeTests(t, cases)
}

func TestBasicMaterialize(t *testing.T) {
	aTime1 := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aDuration := time.Nanosecond * 2785814474
	aDuration2 := time.Nanosecond * 2468348254

	cases := []MaterializeTest{
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: span.Span{
					StartTime: aTime1,
				},
			},
			ExpectedVibe: &Basic{
				Span: span.Span{
					StartTime: aTime1,
				},
				Effects: []ifaces.Effect{
					&effect.Future{
						TimePerLight: &aDuration,
						Painter: &painter.Move{
							ColorStart: &color.WarmCyan,
							Shifter: &shifter.Linear{
								Start:           &aTime1,
								TimePerOneShift: &aDuration2,
							},
						},
					},
				},
			},
		},
	}

	RunMaterializeTests(t, cases)
}
func TestBasicGetStabilizeFuncs(t *testing.T) {
	c := helper.StabilizeableTest{
		Stabalizable: &Basic{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Basic{
				Effects: []ifaces.Effect{
					&effect.Solid{},
				},
			},
			&Basic{
				Effects: []ifaces.Effect{
					&effect.Solid{
						Painter: &painter.Static{},
					},
				},
			},
			&Basic{
				Effects: []ifaces.Effect{
					&effect.Solid{
						Painter: &painter.Static{
							Color: &color.Blue,
						},
					},
				},
			},
		},
		Palette: helper.TestPalette{
			Color:   color.Blue,
			Painter: &painter.Static{},
			Effect:  &effect.Future{},
		},
	}
	helper.RunStabilizeableTest(t, c)
}
