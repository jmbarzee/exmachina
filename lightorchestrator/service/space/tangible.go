package space

type Tangible interface {
	// GetLocation returns the physical location of the device
	GetLocation() Vector
	// SetLocation changes the physical location of the device
	SetLocation(Vector)

	// GetOrientation returns the physical orientation of the device
	GetOrientation() Orientation
	// SetOrientation changes the physical orientation of the device
	SetOrientation(Orientation)

	// GetRotation returns the physical rotation of the device
	GetRotation() Orientation
	// SetRotationchanges the physical rotation of the device
	SetRotation(Orientation)
}
