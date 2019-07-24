package services

import (
	"context"
	"fmt"

	pbd "github.com/jmbarzee/domain/grpc"
	"google.golang.org/grpc"
)

type (
	Service struct {
		ServiceName string
		Port        int
		DomainPort  int
	}
)

func (s Service) Locate(ctx context.Context, serviceName string) ([]string, error) {
	addr := fmt.Sprintf("%127.0.0.1:%v", s.DomainPort)

	domainConn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	serviceRequest := &pbd.GetServicesRequest{
		Name: serviceName,
	}

	domainClient := pbd.NewDomainClient(domainConn)
	reply, err := domainClient.GetServices(ctx, serviceRequest)
	if err != nil {
		return nil, err
	}

	lightOrchastratorAddrs := reply.GetAddresses()
	if len(lightOrchastratorAddrs) == 0 {
		return nil, fmt.Errorf("No address found for %s", serviceName)
	}
	if len(lightOrchastratorAddrs) > 1 {
		fmt.Printf("More than one address found for %s", serviceName)
	}
	return lightOrchastratorAddrs, nil
}
