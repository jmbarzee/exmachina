package device

import (
	"github.com/google/uuid"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
)

// Group represents a group of devices who's effects will share traits
type Group struct {
	BasicDevice
	DeviceNodes []DeviceNode
}

// NewGroup creates a new Group with a unique ID
func NewGroup() Group {
	return Group{
		BasicDevice: BasicDevice{
			ID: uuid.New().String(),
		},
		DeviceNodes: []DeviceNode{},
	}
}

// Allocate passes Vibe into this device and its children
// Allocate stabalize the Vibe before passing it to children devices
func (d Group) Allocate(vibe vibe.Vibe) {
	newVibe := vibe.Stabalize()

	for _, device := range d.DeviceNodes {
		device.Allocate(newVibe)
	}
}

// Insert will attempt to place insert a node into a group until successful
func (d *Group) Insert(parentID string, newNode DeviceNode) error {
	if parentID == d.ID {
		d.DeviceNodes = append(d.DeviceNodes, newNode)
		return nil
	}
	for _, node := range d.DeviceNodes {
		if node.Insert(parentID, newNode) == nil {
			return nil
		}
	}
	return DeviceNodeInsertError
}

// GetChildren returns all groups under the GroupOption
func (d Group) GetChildren() []DeviceNode {
	return d.DeviceNodes
}

// GetType returns the type
func (d Group) GetType() string {
	return "Group"
}
