package domain

import (
	"fmt"
	"sync"
)

// peerMap is a wrapper for sync.map to ensure better type safety
type peerMap struct {
	sMap *sync.Map
}

func (m *peerMap) Delete(uuid string) {
	m.sMap.Delete(uuid)
}

func (m *peerMap) Load(uuid string) (*peer, bool) {
	v, ok := m.sMap.Load(uuid)
	if v == nil {
		return nil, ok
	}
	return v.(*peer), ok
}

func (m *peerMap) LoadOrStore(uuid string, mem *peer) (*peer, bool) {
	v, loaded := m.sMap.LoadOrStore(uuid, mem)
	if !loaded {
		return nil, loaded
	}
	return v.(*peer), loaded
}

func (m *peerMap) Range(f func(uuid string, mem *peer) bool) {
	m.sMap.Range(func(k, v interface{}) bool {
		uuid := k.(string)
		mem := v.(*peer)
		return f(uuid, mem)
	})
}

func (m *peerMap) Store(uuid string, mem *peer) {
	if mem == nil {
		panic(fmt.Errorf("Store() mem was nil"))
	}
	m.sMap.Store(uuid, mem)
}
