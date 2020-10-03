package service

import (
	"sync"

	pb "github.com/jmbarzee/dominion/services/lightorchastrator/grpc"
	device "github.com/jmbarzee/dominion/services/lightorchastrator/service/device"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/pbconvert"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/vibe"
)

// DeviceNodeTree thread-safe tree of allocaters
type DeviceNodeTree struct {
	// RWMutex gates changes to the tree
	rwmutex *sync.RWMutex
	// root is the root allocater
	root device.DeviceNode
}

func (t DeviceNodeTree) Allocate(vibe vibe.Vibe) {
	t.rwmutex.Lock()
	t.root.Allocate(vibe)
	t.rwmutex.Unlock()
}

func (t DeviceNodeTree) Insert(parentID string, newNode device.DeviceNode) error {
	t.rwmutex.Lock()
	err := t.root.Insert(parentID, newNode)
	t.rwmutex.Unlock()
	return err
}

func (t DeviceNodeTree) ToPBDeviceNode() *pb.DeviceNode {
	t.rwmutex.RLock()
	node := pbconvert.NewPBDeviceNode(t.root)
	t.rwmutex.RUnlock()
	return node
}
