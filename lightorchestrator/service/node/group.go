package node

import (
	"errors"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
)

// Group represents a group of devices who's effects will share traits
type Group struct {
	Basic
	Children []Node
}

var _ Node = (*Group)(nil)

// NewGroup creates a new Group with a unique ID
func NewGroup(deviceNodes ...Node) *Group {
	if deviceNodes == nil {
		deviceNodes = []Node{}
	}
	return &Group{
		Basic:    NewBasic(),
		Children: deviceNodes,
	}
}

// Allocate passes Vibe into this device and its children
// Allocate Stabilize the Vibe before passing it to children devices
func (n Group) Allocate(vibe ifaces.Vibe) {
	newVibe := vibe.Stabilize()

	for _, child := range n.Children {
		child.Allocate(newVibe)
		newVibe = newVibe.Duplicate()
	}
}

// Clean removes all stored resources which ended before a time
func (n Group) Clean(t time.Time) {
	for _, child := range n.Children {
		child.Clean(t)
	}
}

// GetChildren returns all groups under the GroupOption
func (n Group) GetChildren() []Node {
	return n.Children
}

// Insert will insert a node underneath a parent node.
func (n *Group) Insert(parentID string, newNode Node) error {
	if parentID == n.ID {
		n.Children = append(n.Children, newNode)
		return nil
	}
	for _, child := range n.Children {
		err := child.Insert(parentID, newNode)
		if err == nil {
			return nil
		} else if errors.Is(err, ParentCantHaveChildrenError) {
			return err
		}
	}
	return FindParentNodeError
}

// Delete will delete a node underneath a parent node.
func (n *Group) Delete(parentID, childID string) error {
	if parentID == n.ID {
		for i, child := range n.Children {
			if child.GetID() == childID {
				length := len(n.Children)
				n.Children[i] = n.Children[length-1] // Copy last element to index i.
				n.Children[length-1] = nil           // Erase last element (write zero value).
				n.Children = n.Children[:length-1]
			}
		}
		return FindChildNodeError
	}
	for _, child := range n.Children {
		err := child.Delete(parentID, childID)
		if err == nil {
			return nil
		} else if errors.Is(err, ParentCantHaveChildrenError) {
			return err
		}
	}
	return FindParentNodeError
}

// GetType returns the type
func (Group) GetType() string {
	return "Group"
}
