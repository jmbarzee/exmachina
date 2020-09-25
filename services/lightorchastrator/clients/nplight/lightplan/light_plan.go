package lightplan

import (
	"container/heap"
	"sync"
	"time"
)

type (
	// LightPlan is an set of light changes which are ordered by their DisplayTime
	LightPlan struct {
		*sync.Mutex
		heap *ChangeHeap
	}
)

// NewLightPlan initializes and returns a LightPlan
func NewLightPlan() LightPlan {
	changeHeap := ChangeHeap(make([]LightChange, 0))
	heap.Init(&changeHeap)
	return LightPlan{
		Mutex: &sync.Mutex{},
		heap:  &changeHeap,
	}
}

// Add inserts a LightChange into a LightPlan
func (p LightPlan) Add(c LightChange) {
	p.Lock()
	heap.Push(p.heap, c)
	p.Unlock()
}

// Advance drops all lightChanges before t and returns the most recent
func (p LightPlan) Advance(t time.Time) *LightChange {
	var change *LightChange
	p.Lock()

	if p.heap.Len() > 0 {
		var past LightChange
		var next LightChange
		next = p.heap.Peek()

		for next.Time.Before(t) {
			past = next
			heap.Pop(p.heap)
			change = &past

			if p.heap.Len() > 0 {
				next = p.heap.Peek()
			} else {
				break
			}
		}
	}

	p.Unlock()
	return change
}
