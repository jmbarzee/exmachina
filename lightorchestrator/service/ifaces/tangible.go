package ifaces

import "github.com/jmbarzee/space"

var _ Tangible = (*space.Object)(nil)

// Tangible is the interface to something exists in space
type Tangible interface {
	// GetLocation returns the physical location of the device
	GetLocation() space.Cartesian
	// SetLocation changes the physical location of the device
	SetLocation(space.Cartesian)

	// GetOrientation returns the physical orientation of the device
	GetOrientation() space.Spherical
	// SetOrientation changes the physical orientation of the device
	SetOrientation(space.Spherical)

	// GetRotation returns the physical rotation of the device
	GetRotation() space.Spherical
	// SetRotationchanges the physical rotation of the device
	SetRotation(space.Spherical)
}
