package shifter

import (
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestFieldShift(t *testing.T) {
	aFloat := 1.1
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Field{
				XBender: &bender.Static{
					TheBend: &aFloat,
				},
				YBender: &bender.Static{
					TheBend: &aFloat,
				},
				ZBender: &bender.Static{
					TheBend: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Light: &light.Basic{
						Location: space.Vector{
							X: 1,
							Y: 2,
							Z: 3,
						},
					},
					ExpectedShift: aFloat * 3,
				},
				{
					Light: &light.Basic{
						Location: space.Vector{
							X: -1,
							Y: -2,
							Z: -3,
						},
					},
					ExpectedShift: aFloat * 3,
				},
				{
					Light: &light.Basic{
						Location: space.Vector{
							X: 0,
							Y: 0,
							Z: 0,
						},
					},
					ExpectedShift: aFloat * 3,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestFieldGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Field{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Field{
				XBender: &bender.Static{},
			},
			&Field{
				XBender: &bender.Static{
					TheBend: &aFloat,
				},
			},
			&Field{
				XBender: &bender.Static{
					TheBend: &aFloat,
				},
				YBender: &bender.Static{
					TheBend: &aFloat, // This is a little dirty. The Benders are both/all pointing to the same struct, so TheBend is set with the first bender
				},
			},
			&Field{
				XBender: &bender.Static{
					TheBend: &aFloat,
				},
				YBender: &bender.Static{
					TheBend: &aFloat,
				},
				ZBender: &bender.Static{
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
