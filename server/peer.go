package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// peer represents another legionnaire in the Legion
// all methods of peer assume that the lock is held by the caller
type Peer struct {
	Identity
	// RWMutex locks everypart of a member except for the UUID (which is always read)
	sync.RWMutex
	// conn is the protocol buffer connection to the member
	conn *grpc.ClientConn
}

// newPeer returns a new peer with the passed identity
func newPeer(identity Identity) *Peer {
	return &Peer{
		Identity: identity,
	}
}

// reconnect assumes caller holds lock !!!
func (p *Peer) reconnect(ctx context.Context) error {
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return err
		}
	}

	// reconnect
	conn, err := grpc.DialContext(ctx, p.addr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	// update peer
	p.conn = conn
	p.LastContact = time.Now()
	return nil
}

// addr assumes caller holds lock !!!
func (p *Peer) addr() string {
	return fmt.Sprintf("%s:%v", p.IP.String(), p.Port)
}

func (d *Domain) checkConnection(peer *Peer) error {
	var err error

	d.debugf(debugLocks, "peer.checkConnection() pre-lock(%v)\n", peer.UUID)
	peer.Lock()
	{
		d.debugf(debugLocks, "peer.checkConnection() in-lock(%v)\n", peer.UUID)

		if peer.conn == nil {
			ctx, cancel := context.WithTimeout(context.Background(), d.config.ConnectionConfig.DialTimeout)
			defer cancel()
			err = peer.reconnect(ctx)
		} else {
			state := peer.conn.GetState()
			switch state {
			// TODO consider trying to reconnect
			case connectivity.Idle:
				// connection is waiting for rpcs
				err = nil
			case connectivity.Connecting:
				// TODO figure out how to handle a connection that is pending
				err = fmt.Errorf("connection is still connecting - uuid:%v", peer.UUID)
			case connectivity.Ready:
				err = nil
			case connectivity.TransientFailure:
				// TODO consider using p.conn.WaitForStateChange()
				err = fmt.Errorf("connection is in Transient Failure - uuid:%v", peer.UUID)
			case connectivity.Shutdown:
				err = fmt.Errorf("connection is shutting down - uuid:%v", peer.UUID)
			default:
				err = fmt.Errorf("connection has un recognized state:%v - uuid:%v", state, peer.UUID)
			}
		}
	}
	peer.Unlock()
	d.debugf(debugLocks, "updateLegion() post-lock(%v)\n", peer.UUID)
	return err
}
