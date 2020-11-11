package pbconvert

import (
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/space"
)

// NewPBVector Builds a Vector for grpc requests
func NewPBVector(vector space.Cartesian) *pb.Cartesian {
	return &pb.Cartesian{
		X: float32(vector.X),
		Y: float32(vector.Y),
		Z: float32(vector.Z),
	}
}

// NewVector Builds a Vector from grpc requests
func NewVector(vector *pb.Cartesian) space.Cartesian {
	return space.Cartesian{
		X: float64(vector.X),
		Y: float64(vector.Y),
		Z: float64(vector.Z),
	}
}

// NewPBOrientation Builds an Orientation for grpc requests
func NewPBOrientation(orientation space.Spherical) *pb.Spherical {
	return &pb.Spherical{
		R: float32(orientation.R),
		T: float32(orientation.T),
		P: float32(orientation.P),
	}
}

// NewOrientation Builds an Orientation from grpc requests
func NewOrientation(orientation *pb.Spherical) space.Spherical {
	return space.Spherical{
		R: float64(orientation.R),
		T: float64(orientation.T),
		P: float64(orientation.P),
	}
}
