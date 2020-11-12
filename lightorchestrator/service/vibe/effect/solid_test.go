package effect

import (
	"testing"
	"time"

	"github.com/jmbarzee/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/painter"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/shifter"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestSolidEffect(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second
	aFloat := 1.0
	numLights := 5
	cases := []EffectTest{
		{
			Name: "Solid Effect with Static Painter",
			Effect: &Solid{
				Painter: &painter.Static{
					Color: color.Blue,
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
			Name: "Solid Effect with Moving Painter",
			Effect: &Solid{
				Painter: &painter.Move{
					ColorStart: color.Blue,
					Shifter: &shifter.Temporal{
						Start:    &aTime,
						Interval: &aSecond,
						Bender: &bender.Linear{
							Interval: &aFloat,
						},
					},
				},
			},
			IntialLights: GetLights(numLights, color.Black),
			Instants: []Instant{
				{
					Time:           aTime.Add(time.Second * 0 / 24),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
				{
					Time:           aTime.Add(time.Second * 4 / 24),
					ExpectedLights: GetLights(numLights, color.Magenta),
				},
				{
					Time:           aTime.Add(time.Second * 8 / 24),
					ExpectedLights: GetLights(numLights, color.Red),
				},
				{
					Time:           aTime.Add(time.Second * 12 / 24),
					ExpectedLights: GetLights(numLights, color.Yellow),
				},
				{
					Time:           aTime.Add(time.Second * 16 / 24),
					ExpectedLights: GetLights(numLights, color.Green),
				},
				{
					Time:           aTime.Add(time.Second * 20 / 24),
					ExpectedLights: GetLights(numLights, color.Cyan),
				},
				{
					Time:           aTime.Add(time.Second * 0 / 24),
					ExpectedLights: GetLights(numLights, color.Blue),
				},
			},
		},
	}
	RunEffectTests(t, cases)
}

func TestSolidGetStabilizeFuncs(t *testing.T) {
	c := helper.StabilizeableTest{
		Stabalizable: &Solid{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Solid{
				Painter: &painter.Static{},
			},
			&Solid{
				Painter: &painter.Static{
					Color: color.Blue,
				},
			},
		},
		Palette: helper.TestPalette{
			Color:   color.Blue,
			Painter: &painter.Static{},
		},
	}
	helper.RunStabilizeableTest(t, c)
}
