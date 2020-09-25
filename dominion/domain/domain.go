package domain

import (
	"time"

	"github.com/jmbarzee/dominion/dominion/config"
	"github.com/jmbarzee/dominion/identity"
	"google.golang.org/grpc"
)

// Domain is a representation of a domain service somewhere on the network
// Domain implements system.Connectable
type Domain struct {
	//DomainIdentity holds the identifying information of the domain
	identity.DomainIdentity

	// conn is the protocol buffer connection to the member
	Conn *grpc.ClientConn

	// LastContact is the last time a domain replied to a rpc
	LastContact time.Time
}

// GetAddress returns the target address for the connection
func (d Domain) GetAddress() identity.Address {
	return d.Address
}

// GetConnection returns a current gRCP connection (for checking)
func (d Domain) GetConnection() *grpc.ClientConn {
	return d.Conn
}

// SetConnection replaces the connection of the device (any existing connection will be closed prior to this)
func (d *Domain) SetConnection(newConn *grpc.ClientConn) {
	d.Conn = newConn
}

// GetTimeout returns the timeout for dialing a new connection
func (Domain) GetTimeout() time.Duration {
	return config.GetDominionConfig().DialTimeout
}
