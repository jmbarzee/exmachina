package bender

import (
	"math"
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestExponentialBend(t *testing.T) {
	exponent := 2.4
	coefficient := 1.1
	cases := []BenderTest{
		{
			Name: "Paint Black",
			Bender: &Exponential{
				Exponent:    &exponent,
				Coefficient: &coefficient,
			},
			Instants: []Instant{
				{
					Input:        -2.0,
					ExpectedBend: coefficient * math.Pow(2.0, exponent),
				},
				{
					Input:        -1.0,
					ExpectedBend: coefficient * math.Pow(1.0, exponent),
				},
				{
					Input:        0.0,
					ExpectedBend: coefficient * math.Pow(0.0, exponent),
				},
				{
					Input:        1.0,
					ExpectedBend: coefficient * math.Pow(1.0, exponent),
				},
				{
					Input:        2.0,
					ExpectedBend: coefficient * math.Pow(2.0, exponent),
				},
			},
		},
	}
	RunBenderTests(t, cases)
}

func TestExponentialGetStabilizeFuncs(t *testing.T) {
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Exponential{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Exponential{
				Exponent: &aFloat,
			},
			&Exponential{
				Exponent:    &aFloat,
				Coefficient: &aFloat,
			},
		},
		Palette: helper.TestPalette{
			Shift: aFloat,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
