package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
)

// NewPBDevice Builds a Device for grpc requests
func NewPBDevice(d device.Device) *pb.Device {
	return &pb.Device{
		UUID:        d.GetID(),
		Type:        d.GetType(),
		Location:    NewPBVector(d.GetLocation()),
		Orientation: NewPBOrientation(d.GetOrientation()),
		Rotation:    NewPBOrientation(d.GetRotation()),
		Nodes:       NewPBNodeList(d.GetNodes()),
	}
}
