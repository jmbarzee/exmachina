package dominion

import (
	"time"

	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/config"
	"google.golang.org/grpc"
)

// Dominion is a representation of the Dominion service somewhere on the network
// Dominion implements system.Connectable
type Dominion struct {
	//DominionIdentity holds the identifying information of the dominion
	identity.DominionIdentity

	// conn is the protocol buffer connection to the member
	Conn *grpc.ClientConn
}

// GetAddress returns the target address for the connection
func (d Dominion) GetAddress() identity.Address {
	return d.Address
}

// GetConnection returns a current gRCP connection (for checking)
func (d Dominion) GetConnection() *grpc.ClientConn {
	return d.Conn
}

// SetConnection replaces the connection of the device (any existing connection will be closed prior to this)
func (d *Dominion) SetConnection(newConn *grpc.ClientConn) {
	d.Conn = newConn
}

// GetTimeout returns the timeout for dialing a new connection
func (Dominion) GetTimeout() time.Duration {
	return config.DefaultServiceDialTimeout
}
