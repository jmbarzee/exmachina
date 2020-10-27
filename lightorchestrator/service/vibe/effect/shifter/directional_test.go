package shifter

import (
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestDirectionalShift(t *testing.T) {
	aFloat := 1.1
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Directional{
				PhiBender: &bender.Static{
					TheBend: &aFloat,
				},
				ThetaBender: &bender.Static{
					TheBend: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Light: &light.Basic{
						Orientation: space.Orientation{
							Phi:   1,
							Theta: 2,
						},
					},
					ExpectedShift: aFloat * 2,
				},
				{
					Light: &light.Basic{
						Orientation: space.Orientation{
							Phi:   -1,
							Theta: -2,
						},
					},
					ExpectedShift: aFloat * 2,
				},
				{
					Light: &light.Basic{
						Orientation: space.Orientation{
							Phi:   0,
							Theta: 0,
						},
					},
					ExpectedShift: aFloat * 2,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestDirectionalGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Directional{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Directional{
				PhiBender: &bender.Static{},
			},
			&Directional{
				PhiBender: &bender.Static{
					TheBend: &aFloat, // This is a little dirty. The Benders are both/all pointing to the same struct, so TheBend is set with the first bender
				},
			},
			&Directional{
				PhiBender: &bender.Static{
					TheBend: &aFloat,
				},
				ThetaBender: &bender.Static{
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
