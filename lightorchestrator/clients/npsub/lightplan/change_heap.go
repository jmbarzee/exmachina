package lightplan

import (
	"time"
)

// An ChangeHeap is a min-heap of LightChange.
type (
	ChangeHeap []LightChange

	// LightChange is a single update to a set of NeoPixel Lights
	LightChange struct {
		Lights []uint32
		Time   time.Time
	}
)

// Len implements the golang heap interface
func (h ChangeHeap) Len() int { return len(h) }

// Less implements the golang heap interface
func (h ChangeHeap) Less(i, j int) bool { return h[i].Time.Before(h[j].Time) }

// Swap implements the golang heap interface
func (h ChangeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push implements the golang heap interface
func (h *ChangeHeap) Push(x interface{}) { *h = append(*h, x.(LightChange)) }

// Pop implements the golang heap interface
func (h *ChangeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return &x
}

// Peek returns the next item in the ChangeHeap without removing it
func (h ChangeHeap) Peek() LightChange {
	if len(h) > 0 {
		return h[0]
	}
	return LightChange{}
}
