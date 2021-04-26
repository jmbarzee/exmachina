package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
)

// NewPBDevice Builds a Device for grpc requests
func NewPBDevice(d device.Device) *pb.Device {
	id := d.GetID()
	return &pb.Device{
		ID:          id[:],
		Type:        d.GetType(),
		Location:    NewPBVector(d.GetLocation()),
		Orientation: NewPBOrientation(d.GetOrientation()),
		Rotation:    NewPBOrientation(d.GetRotation()),
		Nodes:       NewPBNodeList(d.GetNodes()),
	}
}
