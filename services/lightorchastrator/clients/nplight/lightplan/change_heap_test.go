package lightplan

import (
	"container/heap"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestChangeHeap(t *testing.T) {
	atime := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	lightChangeMin := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   atime.Add(time.Hour * -1),
	}
	lightChangeMid := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   atime,
	}
	lightChangeMax := LightChange{
		Lights: []uint32{0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000},
		Time:   atime.Add(time.Hour),
	}

	cases := []struct {
		Name       string
		Initial    ChangeHeap
		Operations func(h *ChangeHeap) *ChangeHeap
		Expected   []LightChange
	}{
		{
			Name:    "Template",
			Initial: ChangeHeap([]LightChange{}),
			Operations: func(h *ChangeHeap) *ChangeHeap {
				return h
			},
			Expected: []LightChange{},
		},
		{
			Name:    "Empty push",
			Initial: ChangeHeap([]LightChange{}),
			Operations: func(h *ChangeHeap) *ChangeHeap {
				heap.Push(h, lightChangeMin)
				return h
			},
			Expected: []LightChange{
				lightChangeMin,
			},
		},
		{
			Name:    "Ordered push",
			Initial: ChangeHeap([]LightChange{}),
			Operations: func(h *ChangeHeap) *ChangeHeap {
				heap.Push(h, lightChangeMin)
				heap.Push(h, lightChangeMid)
				heap.Push(h, lightChangeMax)
				return h
			},
			Expected: []LightChange{
				lightChangeMin,
				lightChangeMid,
				lightChangeMax,
			},
		},
		{
			Name:    "Unordered push",
			Initial: ChangeHeap([]LightChange{}),
			Operations: func(h *ChangeHeap) *ChangeHeap {
				heap.Push(h, lightChangeMax)
				heap.Push(h, lightChangeMid)
				heap.Push(h, lightChangeMin)
				return h
			},
			Expected: []LightChange{
				lightChangeMin,
				lightChangeMid,
				lightChangeMax,
			},
		},
		{
			Name: "Unordered push on filled heap",
			Initial: ChangeHeap([]LightChange{
				lightChangeMin,
				lightChangeMid,
				lightChangeMax}),
			Operations: func(h *ChangeHeap) *ChangeHeap {
				heap.Push(h, lightChangeMax)
				heap.Push(h, lightChangeMid)
				heap.Push(h, lightChangeMin)
				return h
			},
			Expected: []LightChange{
				lightChangeMin,
				lightChangeMin,
				lightChangeMid,
				lightChangeMid,
				lightChangeMax,
				lightChangeMax,
			},
		},
	}
	for i, c := range cases {
		t.Run("case_"+strconv.Itoa(i)+"_"+c.Name, func(t *testing.T) {
			actual := c.Operations(&c.Initial)
			if len(c.Expected) != len(*actual) {
				t.Fatalf("Heap lengths did not match\nExpect %v\nActual %v", len(c.Expected), len(*actual))
			}
			for _, expectedChange := range c.Expected {
				peekedChange := c.Initial.Peek()
				if !reflect.DeepEqual(expectedChange, peekedChange) {
					t.Fatalf("Heap did not pop the expected LightChange\nExpect %v\nPeeked %v", expectedChange, peekedChange)
				}
				actualChange := heap.Pop(&c.Initial)
				if !reflect.DeepEqual(&expectedChange, actualChange) {
					t.Fatalf("Heap did not pop the expected LightChange\nExpect %v\nActual %v", &expectedChange, actualChange)
				}

			}
		})
	}
}
