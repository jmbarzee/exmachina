package pbconvert

import (
	pb "github.com/jmbarzee/dominion/services/lightorchastrator/grpc"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/device"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/device/neopixel"
)

func NewPBDeviceNode(node device.DeviceNode) *pb.DeviceNode {
	pbNode := &pb.DeviceNode{
		UUID: node.GetID(),
	}
	switch n := node.(type) {
	case *device.Group:
		pbNode.Type = "Group"
		for _, child := range n.DeviceNodes {
			pbNode.Children = append(pbNode.Children, NewPBDeviceNode(child))
		}
	case device.GroupOption:
		pbNode.Type = "GroupOption"
		for _, child := range n.Groups {
			pbNode.Children = append(pbNode.Children, NewPBDeviceNode(child))
		}
	case neopixel.Bar:
		pbNode.Type = "npBar"
	}
	return pbNode
}
