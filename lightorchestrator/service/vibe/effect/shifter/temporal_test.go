package shifter

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/vibe/effect/bender"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

func TestTemporalShift(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aFloat := 1.1
	cases := []ShiftTest{
		{
			Name: "One shift per second",
			Shifter: &Temporal{
				Start: &aTime,
				Bender: &bender.Static{
					TheBend: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(0 * time.Second),
					ExpectedShift: aFloat,
				},
				{
					Time:          aTime.Add(1 * time.Second),
					ExpectedShift: aFloat,
				},
				{
					Time:          aTime.Add(1 * time.Hour),
					ExpectedShift: aFloat,
				},
			},
		},
		{
			Name: "One shift per second",
			Shifter: &Temporal{
				Start: &aTime,
				Bender: &bender.Linear{
					Rate: &aFloat,
				},
			},
			Instants: []Instant{
				{
					Time:          aTime.Add(0 * time.Second),
					ExpectedShift: aFloat * 0,
				},
				{
					Time:          aTime.Add(1 * time.Second),
					ExpectedShift: aFloat * 1,
				},
				{
					Time:          aTime.Add(1 * time.Hour),
					ExpectedShift: aFloat * 3600,
				},
			},
		},
	}
	RunShifterTests(t, cases)
}
func TestTemporalGetStabilizeFuncs(t *testing.T) {
	aTime := time.Date(2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	aFloat := 1.1
	c := helper.StabilizeableTest{
		Stabalizable: &Temporal{},
		ExpectedVersions: []ifaces.Stabalizable{
			&Temporal{
				Start: &aTime,
			},
			&Temporal{
				Start:  &aTime,
				Bender: &bender.Static{},
			},
			&Temporal{
				Start: &aTime,
				Bender: &bender.Static{
					TheBend: &aFloat,
				},
			},
		},
		Palette: helper.TestPalette{
			Span: span.Span{
				StartTime: aTime,
			},
			Bender: &bender.Static{},
			Shift:  aFloat,
		},
	}
	helper.RunStabilizeableTest(t, c)
}
