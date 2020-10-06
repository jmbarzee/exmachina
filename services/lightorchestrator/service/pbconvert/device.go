package pbconvert

import (
	pb "github.com/jmbarzee/dominion/services/lightorchestrator/grpc"
	"github.com/jmbarzee/dominion/services/lightorchestrator/service/device"
	"github.com/jmbarzee/dominion/services/lightorchestrator/service/space"
)

func NewPBDevice(device device.Device) *pb.Device {
	return &pb.Device{
		UUID:        device.GetID(),
		Type:        device.GetType(),
		Location:    NewPBVector(device.GetLocation()),
		Orientation: NewPBOrientation(device.GetOrientation()),
	}
}

func NewPBVector(vector space.Vector) *pb.Vector {
	return &pb.Vector{
		X: vector.X,
		Y: vector.Y,
		Z: vector.Z,
	}
}

func NewVector(vector *pb.Vector) space.Vector {
	return space.Vector{
		X: vector.X,
		Y: vector.Y,
		Z: vector.Z,
	}
}

func NewPBOrientation(orientation space.Orientation) *pb.Orientation {
	return &pb.Orientation{
		Theta: orientation.Theta,
		Phi:   orientation.Phi,
	}
}
func NewOrientation(orientation *pb.Orientation) space.Orientation {
	return space.Orientation{
		Theta: orientation.Theta,
		Phi:   orientation.Phi,
	}
}
