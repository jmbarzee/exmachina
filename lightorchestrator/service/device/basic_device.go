package device

import (
	"errors"
)

var DeviceNodeInsertError = errors.New("Failed to find location to insert DeviceNode")

type BasicDevice struct {
	ID string
}

func (d BasicDevice) GetID() string {
	return d.ID
}

func (d BasicDevice) Insert(parentID string, newNode DeviceNode) error {
	return DeviceNodeInsertError
}
