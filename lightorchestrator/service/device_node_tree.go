package service

import (
	"sync"

	pb "github.com/jmbarzee/services/lightorchestrator/grpc"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/pbconvert"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
)

// DeviceNodeTree thread-safe tree of allocaters
type DeviceNodeTree struct {
	// RWMutex gates changes to the tree
	rwmutex *sync.RWMutex
	// root is the root allocater
	root device.DeviceNode
}

// Allocate passes a vibe into the tree where it will be allocated to sub devices as it is stabalized
func (t DeviceNodeTree) Allocate(vibe vibe.Vibe) {
	t.rwmutex.Lock()
	t.root.Allocate(vibe)
	t.rwmutex.Unlock()
}

// Insert places a device in the tree underneath the device with parentID
func (t DeviceNodeTree) Insert(parentID string, newNode device.DeviceNode) error {
	t.rwmutex.Lock()
	err := t.root.Insert(parentID, newNode)
	t.rwmutex.Unlock()
	return err
}

// ToPBDeviceNode converts the nodes in the tree to pb.DeviceNode
func (t DeviceNodeTree) ToPBDeviceNode() *pb.DeviceNode {
	t.rwmutex.RLock()
	node := pbconvert.NewPBDeviceNode(t.root)
	t.rwmutex.RUnlock()
	return node
}
