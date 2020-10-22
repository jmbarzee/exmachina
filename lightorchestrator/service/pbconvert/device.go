package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
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
		X: float32(vector.X),
		Y: float32(vector.Y),
		Z: float32(vector.Z),
	}
}

func NewVector(vector *pb.Vector) space.Vector {
	return space.Vector{
		X: float64(vector.X),
		Y: float64(vector.Y),
		Z: float64(vector.Z),
	}
}

func NewPBOrientation(orientation space.Orientation) *pb.Orientation {
	return &pb.Orientation{
		Theta: float32(orientation.Theta),
		Phi:   float32(orientation.Phi),
	}
}
func NewOrientation(orientation *pb.Orientation) space.Orientation {
	return space.Orientation{
		Theta: float64(orientation.Theta),
		Phi:   float64(orientation.Phi),
	}
}
