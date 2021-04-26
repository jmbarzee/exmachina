package service

import (
	"sync"

	"github.com/google/uuid"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/pbconvert"
)

// NodeTree thread-safe tree of allocaters
type NodeTree struct {
	// RWMutex gates changes to the tree
	rwmutex *sync.RWMutex
	// root is the root allocater
	root node.Node
}

// Allocate passes a vibe into the tree where it will be allocated to sub devices as it is Stabilized
func (t NodeTree) Allocate(vibe ifaces.Vibe) {
	t.rwmutex.Lock()
	t.root.Allocate(vibe)
	t.rwmutex.Unlock()
}

// Insert places a device into the tree underneath the device with parentID
func (t NodeTree) Insert(parentID uuid.UUID, newNode node.Node) error {
	t.rwmutex.Lock()
	err := t.root.Insert(parentID, newNode)
	t.rwmutex.Unlock()
	return err
}

// Delete removes a device from the tree underneath the device with parentID
func (t NodeTree) Delete(parentID, childID uuid.UUID) error {
	t.rwmutex.Lock()
	err := t.root.Delete(parentID, childID)
	t.rwmutex.Unlock()
	return err
}

// ToPBDeviceNode converts the nodes in the tree to pb.DeviceNode
func (t NodeTree) ToPBDeviceNode() *pb.Node {
	t.rwmutex.RLock()
	node := pbconvert.NewPBNode(t.root)
	t.rwmutex.RUnlock()
	return node
}
