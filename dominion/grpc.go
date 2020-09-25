package dominion

import (
	"context"
	"fmt"
	"time"

	"github.com/jmbarzee/dominion/dominion/domain"
	grpc "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/dominion/system/connect"
)

// GetServices implements grpc and allows the domains to use grpc.
// GetServices serves as the directory of services hosted on all domains.
// GetServices is called by services hosted on a single domain to find their dependencies.
func (d *Dominion) GetServices(ctx context.Context, request *grpc.GetServicesRequest) (*grpc.GetServicesReply, error) {
	rpcName := "GetServices"
	system.LogRPCf(rpcName, "Receving request")
	reply := &grpc.GetServicesReply{
		Services: identity.NewPBServiceIdentityList(d.findService(request.Type)),
	}
	system.LogRPCf(rpcName, "Sending reply")
	return reply, nil
}

func (d *Dominion) rpcHeartbeat(ctx context.Context, domainGuard *domain.DomainGuard) {
	rpcName := "Heartbeat"
	uuid := ""
	err := domainGuard.LatchWrite(func(domain *domain.Domain) error {
		uuid = domain.DomainIdentity.UUID

		if err := connect.CheckConnection(ctx, domain); err != nil {
			return fmt.Errorf("Failed to check connection: %w", err)
		}

		// Prepare request
		request := &grpc.HeartbeatRequest{
			Dominion: identity.NewPBDominionIdentity(d.DominionIdentity),
		}

		// Send RPC
		system.LogRPCf(rpcName, "Sending request")
		client := grpc.NewDomainClient(domain.Conn)
		reply, err := client.Heartbeat(ctx, request)
		if err != nil {
			return err
		}
		system.LogRPCf(rpcName, "Recieved reply")

		// Update domain
		domain.LastContact = time.Now()
		fmt.Println(identity.NewDomainIdentity(reply.GetDomain()))
		domain.DomainIdentity = identity.NewDomainIdentity(reply.GetDomain())
		return nil
	})

	if err != nil {
		d.domains.Delete(uuid)
	}
}

func (d *Dominion) rpcStartService(ctx context.Context, domainGuard *domain.DomainGuard, serviceType string) error {
	rpcName := "StartService"
	return domainGuard.LatchWrite(func(domain *domain.Domain) error {

		if err := connect.CheckConnection(ctx, domain); err != nil {
			return fmt.Errorf("Failed to check connection: %w", err)
		}

		// Prepare request
		request := &grpc.StartServiceRequest{
			Type: serviceType,
		}

		// Send RPC
		system.LogRPCf(rpcName, "Sending request")
		client := grpc.NewDomainClient(domain.Conn)
		reply, err := client.StartService(ctx, request)
		if err != nil {
			return err
		}
		system.LogRPCf(rpcName, "Recieved reply")

		// Update domain
		domain.LastContact = time.Now()
		domain.DomainIdentity.Services[serviceType] = identity.NewServiceIdentity(reply.GetService())
		return nil

	})
}
