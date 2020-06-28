package nplight

import (
	"container/heap"
	"sync"
	"time"
)

type (
	LightPlan struct {
		sync.Mutex
		ChangeHeap *ChangeHeap
	}

	LightChange struct {
		Lights []uint32
		Time   time.Time
	}
)

func NewLightPlan() LightPlan {
	var changeHeap ChangeHeap
	heap.Init(&changeHeap)
	return LightPlan{
		ChangeHeap: &changeHeap,
	}
}

func (p LightPlan) Add(c LightChange) {
	p.Lock()
	{
		heap.Push(p.ChangeHeap, c)
	}
	p.Unlock()
}

func (p LightPlan) Peek() LightChange {
	change := LightChange{}
	p.Lock()
	{
		change = p.ChangeHeap.Peek()
	}
	p.Unlock()
	return change
}

func (p LightPlan) Advance() LightChange {
	change := LightChange{}
	p.Lock()
	{
		p.ChangeHeap.Pop()
		change = p.ChangeHeap.Peek()
	}
	p.Unlock()
	return change
}

// An ChangeHeap is a min-heap of LightChange.
type ChangeHeap []LightChange

func (h ChangeHeap) Len() int           { return len(h) }
func (h ChangeHeap) Less(i, j int) bool { return h[i].Time.Before(h[j].Time) }
func (h ChangeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h ChangeHeap) Peek() LightChange {
	if len(h) > 0 {
		return h[0]
	}
	return LightChange{}
}

func (h *ChangeHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(LightChange))
}

func (h *ChangeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
