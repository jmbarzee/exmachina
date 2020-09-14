package service

import (
	"context"
	"fmt"

	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
	"google.golang.org/grpc"
)

// Heartbeat implements grpc and allows the domain to use grpc.
// Heartbeat serves as the heartbeat from a dominion.
func (s *Service) Heartbeat(ctx context.Context, request *pb.ServiceHeartbeatRequest) (*pb.ServiceHeartbeatReply, error) {
	rpcName := "Heartbeat"
	system.LogRPCf(rpcName, "Receving request")

	// Prepare reply
	reply := &pb.ServiceHeartbeatReply{
		Service: identity.NewPBServiceIdentity(s.ServiceIdentity),
	}
	system.LogRPCf(rpcName, "Sending reply")
	return reply, nil
}

func (s Service) RPCGetServices(ctx context.Context, serviceType string) ([]identity.ServiceIdentity, error) {
	rpcName := "GetServices"

	dominionConn, err := grpc.DialContext(ctx, s.DominionIdentity.Address.String(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	serviceRequest := &pb.GetServicesRequest{
		Type: serviceType,
	}

	system.LogRPCf(rpcName, "Sending request")
	dominionClient := pb.NewDominionClient(dominionConn)
	reply, err := dominionClient.GetServices(ctx, serviceRequest)
	if err != nil {
		return nil, err
	}
	system.LogRPCf(rpcName, "Recieved reply")

	services := identity.NewServiceIdentityList(reply.GetServices())
	if len(services) == 0 {
		return nil, fmt.Errorf("No address found for %s", serviceType)
	}

	return services, nil
}
