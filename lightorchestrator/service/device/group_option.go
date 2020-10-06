package device

import (
	"github.com/google/uuid"
	"github.com/jmbarzee/dominion/services/lightorchestrator/service/shared"
	"github.com/jmbarzee/dominion/services/lightorchestrator/service/vibe"
)

// GroupOption represents a series of groups
type GroupOption struct {
	BasicDevice
	Groups []*Group
}

// NewGroupOption creates a new GroupOption with a unique ID
func NewGroupOption() GroupOption {
	return GroupOption{
		BasicDevice: BasicDevice{
			ID: uuid.New().String(),
		},
		Groups: []*Group{},
	}
}

// Allocate passes Vibe into this device and a single child group
// Allocate stabalize the Vibe before passing it to a child group
func (d GroupOption) Allocate(vibe vibe.Vibe) {
	groupNum := shared.RepeatableOption(vibe.Start(), len(d.Groups))
	d.Groups[groupNum].Allocate(vibe)
}

// Insert will attempt to place insert a node into a group until successful
func (d GroupOption) Insert(parentID string, newNode DeviceNode) error {
	for _, group := range d.Groups {
		if group.Insert(parentID, newNode) == nil {
			return nil
		}
	}
	return DeviceNodeInsertError
}

// GetChildren returns all groups under the GroupOption
func (d GroupOption) GetChildren() []DeviceNode {
	nodes := make([]DeviceNode, len(d.Groups))
	for i, group := range d.Groups {
		nodes[i] = group
	}
	return nodes
}

// GetType returns the type
func (d GroupOption) GetType() string {
	return "GroupOption"
}
