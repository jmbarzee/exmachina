package device

import (
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// AllocGroup represents a group of allocaters who's effects will share traits
type AllocGroup struct {
	Allocaters []Allocater
}

// NewAllocGroup creates a new AllocGroup with a unique ID
func NewAllocGroup(allocaters ...Allocater) *AllocGroup {
	if allocaters == nil {
		allocaters = []Allocater{}
	}
	return &AllocGroup{
		Allocaters: allocaters,
	}
}

// Allocate passes Vibe into this device and its children
// Allocate Stabilize the Vibe before passing it to children devices
func (d AllocGroup) Allocate(vibe ifaces.Vibe) {
	newVibe := vibe.Stabilize()

	for _, device := range d.Allocaters {
		device.Allocate(newVibe)
	}
}
