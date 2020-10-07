package device

import (
	"errors"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
)

var DeviceNodeInsertError = errors.New("Failed to find location to insert DeviceNode")

type (
	// A DeviceNode is a node in the device tree
	// DeviceNodes can reference an object which is also a Device
	// DeviceNodes can also be an abstraction which has a Device as a parent or child
	DeviceNode interface {
		// Allocate passes Vibe into this device and its children
		// Allocate typically stabalize the Vibe before passing it to children devices
		Allocate(vibe.Vibe)

		// Insert will place a node underneath a target node.
		// Insert returns an error if the current DeviceNode is a real Device
		Insert(parentID string, newNode DeviceNode) error

		// GetChildren returns any children under the node
		// DeviceNodes which are also Devices will never return their children
		// despite still sometimes passing
		GetChildren() []DeviceNode

		// GetType returns the type
		GetType() string

		// GetID will return the ID of a device node.
		// DeviceNodes which are children of real devices do not have an ID
		// They will return an empty string and thus their location/children can't be modified
		GetID() string
	}

	// Device represents a physical device with lights
	// A device can also exist in the DeviceTree, accessible by a DeviceNode(s)
	Device interface {
		DeviceNode

		// Render produces lights from the effects stored in a device
		Render(time.Time) []shared.Light

		PruneEffects(time.Time)

		// GetLocation returns the physical location of the device
		GetLocation() space.Vector
		// SetLocation changes the physical location of the device
		SetLocation(space.Vector)

		// GetOrientation returns the physical orientation of the device
		GetOrientation() space.Orientation
		// SetOrientation changes the physical orientation of the device
		SetOrientation(space.Orientation)
	}
)
