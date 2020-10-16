package device

import (
	"github.com/jmbarzee/services/lightorchestrator/service/repeatable"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// AllocGroupOption represents a series of AllocGroups
type AllocGroupOption struct {
	AllocGroups []*AllocGroup
}

// NewAllocGroupOption creates a new AllocGroupOption with a unique ID
func NewAllocGroupOption(allocaterGroups ...*AllocGroup) AllocGroupOption {
	if allocaterGroups == nil {
		allocaterGroups = []*AllocGroup{}
	}
	return AllocGroupOption{
		AllocGroups: allocaterGroups,
	}
}

// Allocate passes Vibe into this device and a single child group
// Allocate Stabilize the Vibe before passing it to a child group
func (d AllocGroupOption) Allocate(vibe ifaces.Vibe) {
	groupNum := repeatable.Option(vibe.Start(), len(d.AllocGroups))
	d.AllocGroups[groupNum].Allocate(vibe)
}
