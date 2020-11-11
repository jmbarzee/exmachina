package shifter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

type (
	ShiftTest struct {
		Name     string
		Shifter  ifaces.Shifter
		Instants []Instant
	}

	Instant struct {
		Time          time.Time
		Light         ifaces.Light
		ExpectedShift float64
	}
)

func RunShifterTests(t *testing.T, cases []ShiftTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			for i, instant := range c.Instants {
				actualShift := c.Shifter.Shift(instant.Time, instant.Light)
				if !helper.FloatsEqual(instant.ExpectedShift, actualShift, helper.MinErrColor) {
					t.Fatalf("instant %v failed:\n\tExpected: %v,\n\tActual: %v", i, instant.ExpectedShift, actualShift)
				}
			}
		})
	}
}
