

// Device represents a physical device with lights
// A device is made up of atleast a single Node
type Device interface {
	// Get
	GetNodes() []Node

	// Render produces lights from the effects stored in a device
	Render(time.Time) []light.Light

	space.Tangible
}
