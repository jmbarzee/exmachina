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
	aSpan1 := span.Span{
		StartTime: aTime1,
		EndTime:   aTime1.Add(time.Hour),
	}
	aSpan2 := span.Span{
		StartTime: aTime2,
		EndTime:   aTime2.Add(time.Hour),
	}
	aSpan3 := span.Span{
		StartTime: aTime3,
		EndTime:   aTime3.Add(time.Hour),
	}

	cases := []StabilizeTest{
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: aSpan1,
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: aSpan1,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect: effect.BasicEffect{Span: aSpan1},
						},
					},
				},
				&Basic{
					Span: aSpan1,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect:  effect.BasicEffect{Span: aSpan1},
							TimePerLight: &aDuration,
						},
					},
				},
				&Basic{
					Span: aSpan1,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect:  effect.BasicEffect{Span: aSpan1},
							TimePerLight: &aDuration,
							Painter:      &painter.Move{},
						},
					},
				},
				&Basic{
					Span: aSpan1,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect:  effect.BasicEffect{Span: aSpan1},
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
				Span: aSpan2,
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: aSpan2,
					Effects: []ifaces.Effect{
						&effect.Solid{
							BasicEffect: effect.BasicEffect{Span: aSpan2},
						},
					},
				},
				&Basic{
					Span: aSpan2,
					Effects: []ifaces.Effect{
						&effect.Solid{
							BasicEffect: effect.BasicEffect{Span: aSpan2},
							Painter:     &painter.Static{},
						},
					},
				},
				&Basic{
					Span: aSpan2,
					Effects: []ifaces.Effect{
						&effect.Solid{
							BasicEffect: effect.BasicEffect{Span: aSpan2},
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
				Span: aSpan3,
			},
			ExpectedVibes: []ifaces.Vibe{
				&Basic{
					Span: aSpan3,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect: effect.BasicEffect{Span: aSpan3},
						},
					},
				},
				&Basic{
					Span: aSpan3,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect: effect.BasicEffect{Span: aSpan3},
							Painter:     &painter.Static{},
						},
					},
				},
				&Basic{
					Span: aSpan3,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect: effect.BasicEffect{Span: aSpan3},
							Painter: &painter.Static{
								Color: &color.WarmGreen,
							},
						},
					},
				},
				&Basic{
					Span: aSpan3,
					Effects: []ifaces.Effect{
						&effect.Future{
							BasicEffect:  effect.BasicEffect{Span: aSpan3},
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
	aFloat := 0.277
	aSpan := span.Span{
		StartTime: aTime1,
		EndTime:   aTime1.Add(time.Hour),
	}

	cases := []MaterializeTest{
		{
			Name: "Basic Vibe",
			ActualVibe: &Basic{
				Span: aSpan,
			},
			ExpectedVibe: &Basic{
				Span: aSpan,
				Effects: []ifaces.Effect{
					&effect.Future{
						BasicEffect:  effect.BasicEffect{Span: aSpan},
						TimePerLight: &aDuration,
						Painter: &painter.Move{
							ColorStart: &color.WarmCyan,
							Shifter: &shifter.Static{
								TheShift: &aFloat,
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
