package node

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// A Node is a node in the tree
// Nodes can reference an object which is also a Device or part of a Device
// Node can also be an abstraction which has a Device as a parent or child
type Node interface {
	// Allocate passes Vibe into this device and its children
	// Allocate typically Stabilize the Vibe before passing it to children devices
	Allocate(ifaces.Vibe)
	// Clean removes all stored resources which ended before a time
	Clean(time.Time)

	// GetChildren returns any children under the node
	GetChildren() []Node
	// Insert will insert a node underneath a parent node.
	Insert(parentID uuid.UUID, newNode Node) error
	// Delete will delete a node underneath a parent node.
	Delete(parentID, childID uuid.UUID) error

	// GetType returns the type
	GetType() string
	// GetID will return the ID of a device node.
	GetID() uuid.UUID
}
