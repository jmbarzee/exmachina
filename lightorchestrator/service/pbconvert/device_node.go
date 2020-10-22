package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
)

func NewPBDeviceNode(node device.DeviceNode) *pb.DeviceNode {
	pbNode := &pb.DeviceNode{
		UUID:     node.GetID(),
		Type:     node.GetType(),
		Children: NewPBDeviceNodeList(node.GetChildren()),
	}
	return pbNode
}

func NewPBDeviceNodeList(nodes []device.DeviceNode) []*pb.DeviceNode {
	pbNodes := make([]*pb.DeviceNode, len(nodes))

	for i, node := range nodes {
		pbNodes[i] = NewPBDeviceNode(node)
	}
	return pbNodes
}
