package space

// An Object is something which exists in space
type Object struct {
	// location is the location of the
	location Vector
	// direction is primary orientation of the object
	direction Orientation
	// rotation is secondary orientation of the object
	rotation Orientation
}

// NewObject creates an object
func NewObject(location Vector, direction, rotation Orientation) Object {
	normalizedRotation := direction.PortionOrtagonal(rotation)
	return Object{
		location:  location,
		direction: direction,
		rotation:  normalizedRotation,
	}
}

// GetLocation returns the physical location of the device
func (o Object) GetLocation() Vector {
	return o.location
}

// SetLocation changes the physical location of the device
func (o *Object) SetLocation(newLocation Vector) {
	o.location = newLocation
}

// GetOrientation returns the physical orientation of the device
func (o Object) GetOrientation() Orientation {
	return o.direction
}

// SetOrientation changes the physical orientation of the device
func (o *Object) SetOrientation(newOrientation Orientation) {
	o.direction = newOrientation
	o.rotation = newOrientation.PortionOrtagonal(o.rotation)
}

// GetRotation returns the physical rotation of the device
func (o Object) GetRotation() Orientation {
	return o.direction
}

// SetRotation changes the physical rotation of the device
func (o *Object) SetRotation(newOrientation Orientation) {
	o.rotation = o.direction.PortionOrtagonal(newOrientation)
}
