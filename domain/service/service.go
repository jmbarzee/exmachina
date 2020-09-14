package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmbarzee/dominion/domain/config"
	"github.com/jmbarzee/dominion/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Service struct {
	//ServiceIdentity holds the identifying information of the service
	identity.ServiceIdentity

	// conn is the protocol buffer connection to the member
	Conn *grpc.ClientConn

	// LastContact is the last time a service replied to a rpc
	LastContact time.Time
}

func (s *Service) CheckConnection(ctx context.Context) error {
	if s.Conn == nil {
		ctxDial, cancel := context.WithTimeout(ctx, config.GetDomainConfig().DialTimeout)
		defer cancel()
		if err := s.makeConnection(ctxDial); err != nil {
			return err
		}
	}

	state := s.Conn.GetState()
	switch state {
	// TODO consider trying to reconnect
	case connectivity.Idle:
		// connection is waiting for rpcs
		return nil
	case connectivity.Connecting:
		// TODO figure out how to handle a connection that is pending
		return fmt.Errorf("connection is still connecting - type:%v", s.Type)
	case connectivity.Ready:
		return nil
	case connectivity.TransientFailure:
		// TODO consider using p.conn.WaitForStateChange()
		return fmt.Errorf("connection is in Transient Failure - type:%v", s.Type)
	case connectivity.Shutdown:
		return fmt.Errorf("connection is shutting down - type:%v", s.Type)
	default:
		return fmt.Errorf("connection has un recognized state:%v - type:%v", state, s.Type)
	}
}

func (s *Service) makeConnection(ctx context.Context) error {
	if s.Conn != nil {
		err := s.Conn.Close()
		if err != nil {
			return err
		}
	}

	// connect
	conn, err := grpc.DialContext(ctx, s.Address.String(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	// update service
	s.Conn = conn
	return nil
}
