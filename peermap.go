package domain

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// peerMap is a wrapper for sync.map to ensure better type safety
type peerMap struct {
	sMap *sync.Map
}

func (m *peerMap) Delete(uuid string) {
	m.sMap.Delete(uuid)
}

func (m *peerMap) Load(uuid string) (*Peer, bool) {
	v, ok := m.sMap.Load(uuid)
	if v == nil {
		return nil, ok
	}
	return v.(*Peer), ok
}

func (m *peerMap) LoadOrStore(uuid string, mem *Peer) (*Peer, bool) {
	v, loaded := m.sMap.LoadOrStore(uuid, mem)
	if !loaded {
		return nil, loaded
	}
	return v.(*Peer), loaded
}

func (m *peerMap) Range(f func(uuid string, mem *Peer) bool) {
	m.sMap.Range(func(k, v interface{}) bool {
		uuid := k.(string)
		mem := v.(*Peer)
		return f(uuid, mem)
	})
}

// SizeEstimte only garuntees that the number of all existing keys
// in some length of time is equal to the result
func (m *peerMap) SizeEstimate() int {
	size := int32(0)
	m.sMap.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&size, 1)
		return true
	})
	return int(size)
}

func (m *peerMap) Store(uuid string, mem *Peer) {
	if mem == nil {
		panic(fmt.Errorf("Store() mem was nil"))
	}
	m.sMap.Store(uuid, mem)
}
