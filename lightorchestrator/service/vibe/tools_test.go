package vibe

import (
	"testing"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

type StabilizeTest struct {
	Name          string
	ActualVibe    ifaces.Vibe
	ExpectedVibes []ifaces.Vibe
}

func RunStabilizeTests(t *testing.T, cases []StabilizeTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actualVibe := c.ActualVibe
			for i, expectedVibe := range c.ExpectedVibes {

				actualVibe = actualVibe.Stabilize()
				if !helper.StructsEqual(actualVibe, expectedVibe) {
					t.Fatalf("Stabilize %v failed. Vibes were not equal:\n\tExpected: %+v,\n\tActual: %+v", i, actualVibe, expectedVibe)
				}
			}
		})
	}
}

type MaterializeTest struct {
	Name         string
	ActualVibe   ifaces.Vibe
	ExpectedVibe ifaces.Vibe
}

func RunMaterializeTests(t *testing.T, cases []MaterializeTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.ActualVibe.Materialize()
			if !helper.StructsEqual(c.ActualVibe, c.ExpectedVibe) {
				t.Fatalf("Materialize failed. Vibes were not equal:\n\tExpected: %+v,\n\tActual: %+v", c.ActualVibe, c.ExpectedVibe)
			}
		})
	}
}
