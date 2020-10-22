package lightplan

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestLightPlan(t *testing.T) {
	time0 := time.Date(
		2009, 11, 17, 20, 34, 50, 651387237, time.UTC)
	time1 := time.Date(
		2009, 11, 17, 20, 34, 51, 651387237, time.UTC)
	time2 := time.Date(
		2009, 11, 17, 20, 34, 52, 651387237, time.UTC)
	lightChangeMin := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   time0,
	}
	lightChangeMid := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   time1,
	}
	lightChangeMax := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   time2,
	}

	cases := []struct {
		Name       string
		Initial    LightPlan
		Operations func(h LightPlan) LightPlan
		Expected   []LightChange
	}{
		{
			Name:    "Template",
			Initial: NewLightPlan(),
			Operations: func(p LightPlan) LightPlan {
				return p
			},
			Expected: []LightChange{},
		},
		{
			Name:    "Single LightChange",
			Initial: NewLightPlan(),
			Operations: func(p LightPlan) LightPlan {
				p.Add(lightChangeMin)
				return p
			},
			Expected: []LightChange{
				lightChangeMin,
			},
		},
		{
			Name:    "Tripple LightChange",
			Initial: NewLightPlan(),
			Operations: func(p LightPlan) LightPlan {
				p.Add(lightChangeMin)
				p.Add(lightChangeMid)
				p.Add(lightChangeMax)
				return p
			},
			Expected: []LightChange{
				lightChangeMin,
				lightChangeMid,
				lightChangeMax,
			},
		},
		{
			Name:    "Tripple LightChange (out of order)",
			Initial: NewLightPlan(),
			Operations: func(p LightPlan) LightPlan {
				p.Add(lightChangeMax)
				p.Add(lightChangeMid)
				p.Add(lightChangeMin)
				return p
			},
			Expected: []LightChange{
				lightChangeMin,
				lightChangeMid,
				lightChangeMax,
			},
		},
		{
			Name:    "Duplicate LightChange",
			Initial: NewLightPlan(),
			Operations: func(p LightPlan) LightPlan {
				p.Add(lightChangeMid)
				p.Add(lightChangeMid)
				return p
			},
			Expected: []LightChange{
				lightChangeMid,
			},
		},
	}
	for i, c := range cases {
		t.Run("case_"+strconv.Itoa(i)+"_"+c.Name, func(t *testing.T) {
			actual := c.Operations(c.Initial)
			if len(c.Expected) > actual.heap.Len() {
				t.Fatalf("LightPlan length was not long enough to match\nExpect %v\nActual %v", len(c.Expected), actual.heap.Len())
			}
			for _, expectedChange := range c.Expected {
				timeAfter := expectedChange.Time.Add(time.Millisecond)
				actualChange := actual.Advance(timeAfter)
				if !reflect.DeepEqual(&expectedChange, actualChange) {
					t.Fatalf("LightPlan did not pop the expected LightChange\nExpect %v\nActual %v", &expectedChange, actualChange)
				}
			}
		})
	}
}
