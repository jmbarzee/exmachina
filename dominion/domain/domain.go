package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/jmbarzee/dominion/dominion/config"
	"github.com/jmbarzee/dominion/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Domain struct {
	//DomainIdentity holds the identifying information of the domain
	identity.DomainIdentity

	// conn is the protocol buffer connection to the member
	Conn *grpc.ClientConn

	// LastContact is the last time a domain replied to a rpc
	LastContact time.Time
}

func (d *Domain) CheckConnection(ctx context.Context) error {
	if d.Conn == nil {
		ctxDial, cancel := context.WithTimeout(ctx, config.GetDominionConfig().DialTimeout)
		defer cancel()
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
		return fmt.Errorf("connection is still connecting - uuid:%v", d.UUID)
	case connectivity.Ready:
		return nil
	case connectivity.TransientFailure:
		// TODO consider using p.conn.WaitForStateChange()
		return fmt.Errorf("connection is in Transient Failure - uuid:%v", d.UUID)
	case connectivity.Shutdown:
		return fmt.Errorf("connection is shutting down - uuid:%v", d.UUID)
	default:
		return fmt.Errorf("connection has un recognized state:%v - uuid:%v", state, d.UUID)
	}
}

func (d *Domain) makeConnection(ctx context.Context) error {
	if d.Conn != nil {
		err := d.Conn.Close()
		if err != nil {
			return err
		}
	}

	// connect
	conn, err := grpc.DialContext(ctx, d.Address.String(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	// update domain
	d.Conn = conn
	return nil
}
