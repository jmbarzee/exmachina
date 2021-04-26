package device

import "github.com/google/uuid"

// Basic implements some traits and features which are shared between all nodes
type Basic struct {
	id uuid.UUID
}

// NewBasic creates a Basic
func NewBasic(id uuid.UUID) Basic {
	return Basic{
		id: id,
	}
}

func (d Basic) GetID() uuid.UUID {
	return d.id
}
