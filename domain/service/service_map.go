package service

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/jmbarzee/dominion/system"
)

// ServiceMap is a wrapper for sync.map to ensure better type safety
type ServiceMap struct {
	sMap *sync.Map
}

func NewServiceMap() ServiceMap {
	return ServiceMap{
		sMap: &sync.Map{},
	}
}

func (m ServiceMap) Delete(uuid string) {
	m.sMap.Delete(uuid)
}

func (m ServiceMap) Load(uuid string) (*ServiceGuard, bool) {
	v, ok := m.sMap.Load(uuid)
	if v == nil {
		return nil, ok
	}
	return v.(*ServiceGuard), ok
}

func (m ServiceMap) LoadOrStore(uuid string, mem *ServiceGuard) (*Service, bool) {
	v, loaded := m.sMap.LoadOrStore(uuid, mem)
	if !loaded {
		return nil, loaded
	}
	return v.(*Service), loaded
}

func (m ServiceMap) Range(f func(uuid string, mem *ServiceGuard) bool) {
	m.sMap.Range(func(k, v interface{}) bool {
		uuid := k.(string)
		mem := v.(*ServiceGuard)
		return f(uuid, mem)
	})
}

// SizeEstimte only garuntees that the number of all existing keys
// in some length of time is equal to the result
func (m *ServiceMap) SizeEstimate() int {
	size := int32(0)
	m.sMap.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&size, 1)
		return true
	})
	return int(size)
}

func (m *ServiceMap) Store(uuid string, mem *ServiceGuard) {
	if mem == nil {
		system.Panic(fmt.Errorf("Store() mem was nil"))
	}
	m.sMap.Store(uuid, mem)
}
