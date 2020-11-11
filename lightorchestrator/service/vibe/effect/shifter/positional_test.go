package shifter

import (
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestPositionalShift(t *testing.T) {
	aPosition := 5
	numPositions := 25
	aFloat := 1.1
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Positional{
				Bender: &bender.Static{
					TheBend: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Light: &light.Basic{
						Position:     aPosition,
						NumPositions: numPositions,
					},
					ExpectedShift: aFloat,
				},
			},
		},
		{
			Name: "One shift per second",
			Shifter: &Positional{
				Bender: &bender.Linear{
					Interval: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Light: &light.Basic{
						Position:     aPosition,
						NumPositions: numPositions,
					},
					ExpectedShift: float64(aPosition) / aFloat / float64(numPositions),
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestPositionalGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Positional{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Positional{
				Bender: &bender.Static{},
			},
			&Positional{
				Bender: &bender.Static{
					TheBend: &aFloat,
				},
			},
		},
		Palette: helper.TestPalette{
			Bender: &bender.Static{},
			Shift:  aFloat,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
