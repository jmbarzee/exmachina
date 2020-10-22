package shifter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestLinearShift(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aSecond := time.Second
	aMinute := time.Minute
	anHour := time.Hour
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Linear{
				Start:           &aTime,
				TimePerOneShift: &aSecond,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Millisecond),
					ExpectedShift: 1.0 / 1000,
				},
				{
					Time:          aTime.Add(time.Second),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Minute),
					ExpectedShift: 60.0,
				},
				{
					Time:          aTime.Add(time.Hour),
					ExpectedShift: 3600.0,
				},
			},
		},
		{
			Name: "One shift per minute",
			Shifter: &Linear{
				Start:           &aTime,
				TimePerOneShift: &aMinute,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Millisecond),
					ExpectedShift: 1.0 / 60 / 1000,
				},
				{
					Time:          aTime.Add(time.Second),
					ExpectedShift: 1.0 / 60,
				},
				{
					Time:          aTime.Add(time.Minute),
					ExpectedShift: 1,
				},
				{
					Time:          aTime.Add(time.Hour),
					ExpectedShift: 60.0,
				},
			},
		},
		{
			Name: "One shift per hour",
			Shifter: &Linear{
				Start:           &aTime,
				TimePerOneShift: &anHour,
			},
			Instants: []Instant{
				{
					Time:          aTime,
					ExpectedShift: 0,
				},
				{
					Time:          aTime.Add(time.Millisecond),
					ExpectedShift: 1.0 / 3600 / 1000,
				},
				{
					Time:          aTime.Add(time.Second),
					ExpectedShift: 1.0 / 3600,
				},
				{
					Time:          aTime.Add(time.Minute),
					ExpectedShift: 1.0 / 60,
				},
				{
					Time:          aTime.Add(time.Hour),
					ExpectedShift: 1.0,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestLinearGetStabilizeFuncs(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aDuration := time.Second
	c := helper.StabilizeableTest{
		Stabalizable: &Linear{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Linear{
				Start: &aTime,
			},
			&Linear{
				Start:           &aTime,
				TimePerOneShift: &aDuration,
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Duration: aDuration,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
