package shifter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestSinusoidalShift(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second
	aMinute := time.Minute
	aOne := 1.0
	aTwo := 2.0
	cases := []ShiftTest{
		{
			Name: "Shift by 1 every second",
			Shifter: &Sinusoidal{
				Start:         &aTime,
				TimePerCycle:  &aSecond,
				ShiftPerCycle: &aOne,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: .5,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 4),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 4),
					ExpectedShift: .5,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 4),
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 4),
					ExpectedShift: .5,
				},
			},
		},
		{
			Name: "Shift by 2 every second",
			Shifter: &Sinusoidal{
				Start:         &aTime,
				TimePerCycle:  &aSecond,
				ShiftPerCycle: &aTwo,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Second * 1 / 4),
					ExpectedShift: 2,
				},
				{
					Time:          aTime.Add(time.Second * 2 / 4),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Second * 3 / 4),
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Second * 4 / 4),
					ExpectedShift: 1,
				},
			},
		},
		{
			Name: "Shift by 1 every minute",
			Shifter: &Sinusoidal{
				Start:         &aTime,
				TimePerCycle:  &aMinute,
				ShiftPerCycle: &aOne,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: .5,
				},
				{
					Time:          aTime.Add(time.Minute * 1 / 4),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Minute * 2 / 4),
					ExpectedShift: .5,
				},
				{
					Time:          aTime.Add(time.Minute * 3 / 4),
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Minute * 4 / 4),
					ExpectedShift: .5,
				},
			},
		},
		{
			Name: "Shift by 2 every minute",
			Shifter: &Sinusoidal{
				Start:         &aTime,
				TimePerCycle:  &aMinute,
				ShiftPerCycle: &aTwo,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Minute * 1 / 4),
					ExpectedShift: 2,
				},
				{
					Time:          aTime.Add(time.Minute * 2 / 4),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Minute * 3 / 4),
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Minute * 4 / 4),
					ExpectedShift: 1,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}

func TestSinusoidalGetStabilizeFuncs(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aDuration := time.Second
	aShift := float64(1.0)
	c := helper.StabilizeableTest{
		Stabalizable: &Sinusoidal{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Sinusoidal{
				Start: &aTime,
			},
			&Sinusoidal{
				Start:        &aTime,
				TimePerCycle: &aDuration,
			},
			&Sinusoidal{
				Start:         &aTime,
				TimePerCycle:  &aDuration,
				ShiftPerCycle: &aShift,
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Duration: aDuration,
			Shift:    aShift,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
