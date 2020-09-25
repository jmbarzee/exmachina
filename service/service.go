package service

import (
	"bufio"
	"fmt"
	"os"

	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/service/dominion"
	"github.com/jmbarzee/dominion/system"
	"google.golang.org/grpc"
)

type (
	// Service offers all the shared features of services
	// Service should be emmbeded into the implementation of a specific service
	// the specific service should implement myService.Run(ctx)
	// and should return myService.HostService() (blocking) as its final line
	Service struct {
		UUID string

		identity.ServiceIdentity

		Server *grpc.Server

		Dominion *dominion.DominionGuard
	}
)

// NewService builds a service from a ServiceConfig
func NewService(config config.ServiceConfig) (*Service, error) {
	if err := system.Setup(config.DomainUUID, config.ServiceType); err != nil {
		return nil, err
	}

	if err := captureStdout(); err != nil {
		return nil, err
	}

	if err := captureStderr(); err != nil {
		return nil, err
	}

	// Initialize IP
	ip, err := system.GetOutboundIP()
	if err != nil {
		return nil, fmt.Errorf("failed to find Local IP: %w", err)
	}

	server := grpc.NewServer()

	service := &Service{
		UUID: config.DomainUUID,
		ServiceIdentity: identity.ServiceIdentity{
			Type: config.ServiceType,
			Address: identity.Address{
				IP:   ip,
				Port: config.ServicePort,
			},
		},
		Server: server,
		Dominion: dominion.NewDominionGuard(identity.DominionIdentity{
			Address: identity.Address{
				IP:   config.DominionIP,
				Port: config.DominionPort,
			},
		}),
	}

	pb.RegisterServiceServer(service.Server, service)
	return service, nil
}

func captureStdout() error {
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to gather: %w", err)
	}
	os.Stdout = w

	go func() {
		routineName := "Nameless-captureStdout"
		system.LogRoutinef(routineName, "Starting routine")
		buf := bufio.NewReader(r)
		var err error
		var bytes []byte
		for err == nil {
			bytes, err = buf.ReadBytes('\n')
			if err != nil {
				system.Logf("Stdout %s", bytes)
			}
		}
		system.Errorf("Failure while reading bytes from Stdout to log")
		system.LogRoutinef(routineName, "Stopping routine")
	}()
	return nil
}

func captureStderr() error {
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to gather: %w", err)
	}
	os.Stderr = w

	go func() {
		routineName := "Nameless-captureStderr"
		system.LogRoutinef(routineName, "Starting routine")
		buf := bufio.NewReader(r)
		var err error
		var bytes []byte
		for err == nil {
			bytes, err = buf.ReadBytes('\n')
			if err != nil {
				system.Logf("Stderr %s", bytes)
			}
		}
		system.Errorf("Failure while reading bytes from Stderr to log")
		system.LogRoutinef(routineName, "Stopping routine")
	}()
	return nil
}
