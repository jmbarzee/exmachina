package device

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
)

// Device represents a physical device with lights
// A device is made up of atleast a single Node
type Device interface {
	// GetNodes returns all the Nodes which the device holds
	GetNodes() []node.Node

	// Render produces lights from the effects stored in a device
	Render(time.Time) []ifaces.Light

	// GetType returns the type
	GetType() string
	// GetID will return the ID of a device node.
	GetID() uuid.UUID

	ifaces.Tangible
}
