package domain

import (
	"sync"

	"github.com/jmbarzee/dominion/identity"
)

// DomainGuard protects a domain for concurrent access
type DomainGuard struct {
	// RWMutex controls gates to the domain's data
	rwmutex sync.RWMutex
	// data is the Domains
	domain Domain
}

// NewDomainGuard returns a new DomainGuard with the passed identity
func NewDomainGuard(identity identity.DomainIdentity) *DomainGuard {
	return &DomainGuard{
		domain: Domain{
			DomainIdentity: identity,
		},
	}
}

func (d *DomainGuard) LatchWrite(operation func(*Domain) error) error {
	d.rwmutex.Lock()
	err := operation(&d.domain)
	d.rwmutex.Unlock()
	return err
}

func (d *DomainGuard) LatchRead(operation func(*Domain) error) error {
	d.rwmutex.RLock()
	err := operation(&d.domain)
	d.rwmutex.RUnlock()
	return err
}
