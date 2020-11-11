package device

// Basic implements some traits and features which are shared between all nodes
type Basic struct {
	ID string
}

// NewBasic creates a Basic
func NewBasic(uuid string) Basic {
	return Basic{
		ID: uuid,
	}
}

func (d Basic) GetID() string {
	return d.ID
}
