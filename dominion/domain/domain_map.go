package domain

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/jmbarzee/dominion/system"
)

// DomainMap is a wrapper for sync.map to ensure better type safety
type DomainMap struct {
	sMap *sync.Map
}

func NewDomainMap() DomainMap {
	return DomainMap{
		sMap: &sync.Map{},
	}
}

func (m DomainMap) Delete(uuid string) {
	m.sMap.Delete(uuid)
}

func (m DomainMap) Load(uuid string) (*DomainGuard, bool) {
	v, ok := m.sMap.Load(uuid)
	if v == nil {
		return nil, ok
	}
	return v.(*DomainGuard), ok
}

func (m DomainMap) LoadOrStore(uuid string, mem *DomainGuard) (*Domain, bool) {
	v, loaded := m.sMap.LoadOrStore(uuid, mem)
	if !loaded {
		return nil, loaded
	}
	return v.(*Domain), loaded
}

func (m DomainMap) Range(f func(uuid string, mem *DomainGuard) bool) {
	m.sMap.Range(func(k, v interface{}) bool {
		uuid := k.(string)
		mem := v.(*DomainGuard)
		return f(uuid, mem)
	})
}

// SizeEstimte only garuntees that the number of all existing keys
// in some length of time is equal to the result
func (m *DomainMap) SizeEstimate() int {
	size := int32(0)
	m.sMap.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&size, 1)
		return true
	})
	return int(size)
}

func (m *DomainMap) Store(uuid string, mem *DomainGuard) {
	if mem == nil {
		system.Panic(fmt.Errorf("Store() mem was nil"))
	}
	m.sMap.Store(uuid, mem)
}
