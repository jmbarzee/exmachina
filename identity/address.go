package identity

import (
	"fmt"
	"net"

	pb "github.com/jmbarzee/dominion/grpc"
)

// Address is a IP and port combination
type Address struct {
	// IP is the ip which the Domain will be responding on
	IP net.IP
	// Port is the port which the Domain will be responding on
	Port int
}

// NewAddress creates a Address from a pb.Address
func NewAddress(pbAddr *pb.Address) Address {
	return Address{
		IP:   net.ParseIP(string(pbAddr.GetIP())),
		Port: int(pbAddr.GetPort()),
	}
}

// NewPBAddress creates a pb.ServiceIdentity from a Address
func NewPBAddress(addr Address) *pb.Address {
	return &pb.Address{
		IP:   []byte(addr.IP.String()),
		Port: int32(addr.Port),
	}
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%v", a.IP.String(), a.Port)
}
