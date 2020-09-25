package device

type BasicDevice struct {
	ID string
}

func (d BasicDevice) GetID() string {
	return d.ID
}

func (d BasicDevice) Insert(parentID string, newNode DeviceNode) error {
	return DeviceNodeInsertError
}
