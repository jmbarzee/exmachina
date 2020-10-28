package effect

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/painter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestFutureEffect(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aFloat := 1.0
	aSecond := time.Second
	a24thSecond := time.Second / 24
	numLights := 5
	cases := []EffectTest{
		{
			Name: "Future Effect with Static Painter",
			Effect: &Future{
				TimePerLight: &aSecond,
				Painter: &painter.Static{
					Color: &color.Blue,
				},
			},
			IntialLights: GetLights(numLights, color.Black),
			Instants: []Instant{
				{
					Time:           aTime,
					ExpectedLights: GetLights(numLights, color.Blue),
				},
				{
					Time:           aTime.Add(time.Millisecond * 1),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
				{
					Time:           aTime.Add(time.Second * 1),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
				{
					Time:           aTime.Add(time.Minute * 1),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
				{
					Time:           aTime.Add(time.Hour * 1),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
			},
		},
		{
			Name: "Future Effect with Moving Painter",
			Effect: &Future{
				TimePerLight: &a24thSecond,
				Painter: &painter.Move{
					ColorStart: &color.Blue,
					Shifter: &shifter.Temporal{
						Start:    &aTime,
						Interval: &aSecond,
						Bender: &bender.Linear{
							Interval: &aFloat,
						},
					},
				},
			},
			IntialLights: GetLights(3, color.Black),
			Instants: []Instant{
				{
					Time: aTime.Add(time.Second * 0 / 24),
					ExpectedLights: []light.Light{
						&TestLight{
							Color: color.Blue,
						},
						&TestLight{
							Color: color.WarmBlue,
						},
						&TestLight{
							Color: color.Violet,
						},
					},
				},
				{
					Time: aTime.Add(time.Second * 1 / 24),
					ExpectedLights: []light.Light{
						&TestLight{
							Color: color.WarmBlue,
						},
						&TestLight{
							Color: color.Violet,
						},
						&TestLight{
							Color: color.CoolMagenta,
						},
					},
				},
				{
					Time: aTime.Add(time.Second * 2 / 24),
					ExpectedLights: []light.Light{
						&TestLight{
							Color: color.Violet,
						},
						&TestLight{
							Color: color.CoolMagenta,
						},
						&TestLight{
							Color: color.Magenta,
						},
					},
				},
			},
		},
	}
	RunEffectTests(t, cases)
}

func TestFutureGetStabilizeFuncs(t *testing.T) {
	aSecond := time.Second
	c := helper.StabilizeableTest{
		Stabalizable: &Future{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Future{
				TimePerLight: &aSecond,
			},
			&Future{
				TimePerLight: &aSecond,
				Painter:      &painter.Static{},
			},
			&Future{
				TimePerLight: &aSecond,
				Painter: &painter.Static{
					Color: &color.Blue,
				},
			},
		},
		Palette: helper.TestPalette{
			Color:    color.Blue,
			Painter:  &painter.Static{},
			Duration: aSecond,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
