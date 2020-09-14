package domain

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmbarzee/dominion/domain/dominion"
	service "github.com/jmbarzee/dominion/domain/service"
	grpc "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
)

// Heartbeat implements grpc and allows the domain to use grpc.
// Heartbeat serves as the heartbeat from a dominion.
func (d *Domain) Heartbeat(ctx context.Context, request *grpc.HeartbeatRequest) (*grpc.HeartbeatReply, error) {
	rpcName := "Heartbeat"
	system.LogRPCf(rpcName, "Receving request")
	if err := d.updateDominion(identity.NewDominionIdentity(request.GetDominion())); err != nil {
		return nil, err
	}

	// Prepare reply
	reply := &grpc.HeartbeatReply{
		Domain: identity.NewPBDomainIdentity(d.packageDomainIdentity()),
	}
	system.LogRPCf(rpcName, "Sending reply")
	return reply, nil
}

func (d *Domain) rpcHeartbeat(ctx context.Context, serviceGuard *service.ServiceGuard) {
	time.Sleep(time.Second * 3)
	rpcName := "Heartbeat"
	serviceType := ""
	err := serviceGuard.LatchWrite(func(service *service.Service) error {
		serviceType = service.Type

		if err := service.CheckConnection(ctx); err != nil {
			return fmt.Errorf("Failed to check connection: %w", err)
		}

		// Prepare request
		request := &grpc.ServiceHeartbeatRequest{
			Domain: identity.NewPBDomainIdentity(d.DomainIdentity),
		}

		// Send RPC
		system.LogRPCf(rpcName, "Sending request")
		client := grpc.NewServiceClient(service.Conn)
		reply, err := client.Heartbeat(ctx, request)
		if err != nil {
			return err
		}
		system.LogRPCf(rpcName, "Recieved reply")

		// Update domain
		service.LastContact = time.Now()
		service.ServiceIdentity = identity.NewServiceIdentity(reply.GetService())
		return nil
	})

	if err != nil {
		system.Logf("Failed to heartbeat \"%v\": %v: Dropping service", serviceType, err.Error())
		d.services.Delete(serviceType)
	}
}

// StartService implements grpc and initiates a service in the domain.
func (d *Domain) StartService(ctx context.Context, request *grpc.StartServiceRequest) (*grpc.StartServiceReply, error) {
	rpcName := "StartService"
	system.LogRPCf(rpcName, "Receving request")
	ident, err := d.startService(request.GetType())
	if err != nil {
		return nil, fmt.Errorf("Failed to start service: %w", err)
	}

	reply := &grpc.StartServiceReply{
		Service: identity.NewPBServiceIdentity(ident),
	}

	system.LogRPCf(rpcName, "Sending reply")
	return reply, nil
}

func (d *Domain) startService(serviceType string) (identity.ServiceIdentity, error) {
	if _, ok := d.services.Load(serviceType); ok {
		return identity.ServiceIdentity{}, fmt.Errorf("Service already exists! (%s)", serviceType)
	}

	var dominionAddr identity.Address
	d.Dominion.LatchRead(func(dominion *dominion.Dominion) error {
		dominionAddr = dominion.Address
		return nil
	})
	dominionIP := dominionAddr.IP
	dominionPort := dominionAddr.Port
	servicePort := rand.Intn((1 << 16) - 1)

	err := service.Start(serviceType, dominionIP, dominionPort, servicePort)
	if err != nil {
		return identity.ServiceIdentity{}, err
	}

	serviceIdent := identity.ServiceIdentity{
		Type: serviceType,
		Address: identity.Address{
			IP:   d.Address.IP,
			Port: servicePort,
		},
	}

	d.services.Store(serviceType, service.NewServiceGuard(serviceIdent))

	system.Logf("Started service: \"%v\" at port:%v", serviceType, servicePort)

	return serviceIdent, nil
}
