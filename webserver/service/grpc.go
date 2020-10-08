package service

import (
	"context"

	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/dominion"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/dominion/system/connect"
)

func (s WebServer) rpcGetDomains(ctx context.Context) ([]identity.DomainIdentity, error) {
	rpcName := "GetDomain"
	domains := []identity.DomainIdentity{}

	err := s.Dominion.LatchWrite(func(dominion *dominion.Dominion) error {
		err := connect.CheckConnection(ctx, dominion)
		if err != nil {
			return err
		}

		serviceRequest := &pb.Empty{}

		system.LogRPCf(rpcName, "Sending request")
		dominionClient := pb.NewDominionClient(dominion.Conn)
		reply, err := dominionClient.GetDomains(ctx, serviceRequest)
		if err != nil {
			return err
		}
		system.LogRPCf(rpcName, "Received reply")

		domains = identity.NewDomainIdentityList(reply.GetDomains())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return domains, nil
}
