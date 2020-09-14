package dominion

import (
	"context"
	"fmt"
	"time"

	"github.com/jmbarzee/dominion/dominion/config"
	"github.com/jmbarzee/dominion/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Dominion struct {
	//DominionIdentity holds the identifying information of the service
	identity.DominionIdentity

	// conn is the protocol buffer connection to the member
	Conn *grpc.ClientConn

	// LastContact is the last time a service replied to a rpc
	LastContact time.Time
}

func (d *Dominion) CheckConnection(ctx context.Context) error {
	if d.Conn == nil {
		ctxDial, _ := context.WithTimeout(ctx, config.GetDominionConfig().DialTimeout)
		if err := d.makeConnection(ctxDial); err != nil {
			return err
		}
	}

	state := d.Conn.GetState()
	switch state {
	// TODO consider trying to reconnect
	case connectivity.Idle:
		// connection is waiting for rpcs
		return nil
	case connectivity.Connecting:
		// TODO figure out how to handle a connection that is pending
		return fmt.Errorf("connection is still connecting")
	case connectivity.Ready:
		return nil
	case connectivity.TransientFailure:
		// TODO consider using p.conn.WaitForStateChange()
		return fmt.Errorf("connection is in Transient Failure")
	case connectivity.Shutdown:
		return fmt.Errorf("connection is shutting down")
	default:
		return fmt.Errorf("connection has un recognized state:%v", state)
	}
}

func (s *Dominion) makeConnection(ctx context.Context) error {
	if s.Conn != nil {
		err := s.Conn.Close()
		if err != nil {
			return err
		}
	}

	// connect
	conn, err := grpc.DialContext(ctx, s.addr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	// update service
	s.Conn = conn
	s.LastContact = time.Now()
	return nil
}

func (s *Dominion) addr() string {
	return fmt.Sprintf("%s:%v", s.Address.IP.String(), s.Address.Port)
}
