package node

import (
	"errors"

	"github.com/google/uuid"
)

var FindParentNodeError = errors.New("Failed to find parent node with matching ID")
var FindChildNodeError = errors.New("Failed to find child node with matching ID")
var ParentCantHaveChildrenError = errors.New("Found node with matching ID, but node can't hold children")

type Basic struct {
	ID string
}

func NewBasic() Basic {
	return Basic{
		ID: uuid.New().String(),
	}
}

// GetChildren returns any children under the node
func (Basic) GetChildren() []Node {
	return nil
}

// Insert will insert a node underneath a parent node.
func (d Basic) Insert(parentID string, newNode Node) error {
	if parentID == d.ID {
		return ParentCantHaveChildrenError
	}
	return FindParentNodeError
}

// Delete will delete a node underneath a parent node.
func (d Basic) Delete(parentID, childID string) error {
	if parentID == d.ID {
		return ParentCantHaveChildrenError
	}
	return FindParentNodeError
}

func (d Basic) GetID() string {
	return d.ID
}
