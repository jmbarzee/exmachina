package shifter

import (
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestStaticShift(t *testing.T) {
	aFloat := 1.1
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Static{
				TheShift: &aFloat,
			},
			Instants: []Instant{
				{
					ExpectedShift: aFloat,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestStaticGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Static{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Static{
				TheShift: &aFloat,
			},
		},
		Palette: helper.TestPalette{
			Shift: aFloat,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
