package bender

import (
	"math"
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestSinusoidalBend(t *testing.T) {
	offset := 2.4
	period := 1.1
	amplitude := 3.7
	cases := []BenderTest{
		{
			Name: "Paint Black",
			Bender: &Sinusoidal{
				Offset:    &offset,
				Period:    &period,
				Amplitude: &amplitude,
			},
			Instants: []Instant{
				{
					Input:        -2.5,
					ExpectedBend: amplitude * math.Sin(offset+2*math.Pi*(-2.5/period)),
				},
				{
					Input:        -1.0,
					ExpectedBend: amplitude * math.Sin(offset+2*math.Pi*(-1.0/period)),
				},
				{
					Input:        0.0,
					ExpectedBend: amplitude * math.Sin(offset+2*math.Pi*(0.0/period)),
				},
				{
					Input:        1.0,
					ExpectedBend: amplitude * math.Sin(offset+2*math.Pi*(1.0/period)),
				},
				{
					Input:        2.5,
					ExpectedBend: amplitude * math.Sin(offset+2*math.Pi*(2.5/period)),
				},
			},
		},
	}
	RunBenderTests(t, cases)
}

func TestSinusoidalGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Sinusoidal{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Sinusoidal{
				Offset: &aFloat,
			},
			&Sinusoidal{
				Offset: &aFloat,
				Period: &aFloat,
			},
			&Sinusoidal{
				Offset:    &aFloat,
				Period:    &aFloat,
				Amplitude: &aFloat,
			},
		},
		Palette: helper.TestPalette{
			Shift: aFloat,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
