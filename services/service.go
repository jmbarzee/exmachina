package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	pbd "github.com/jmbarzee/domain/server/grpc"
	"github.com/jmbarzee/domain/server/identity"
	"google.golang.org/grpc"
)

type (
	Service struct {
		ServiceName string
		Port        int
		DomainPort  int
		Logger      *log.Logger
	}
)

func (s Service) DomainAddress() string {
	return fmt.Sprintf("127.0.0.1:%v", s.DomainPort)
}

func (s Service) Locate(ctx context.Context, serviceName string) ([]string, error) {
	s.Logger.Printf("[Service] Locate: %v", serviceName)

	domainConn, err := grpc.DialContext(ctx, s.DomainAddress(), grpc.WithInsecure(), grpc.WithBlock())
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

	addrs := reply.GetAddresses()
	if len(addrs) == 0 {
		return nil, fmt.Errorf("No address found for %s", serviceName)
	}
	return addrs, nil
}

func (s Service) Dump(ctx context.Context) (identity.Identity, []identity.Identity, error) {
	s.Logger.Printf("[Service] Dump")
	domainConn, err := grpc.DialContext(ctx, s.DomainAddress(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return identity.Identity{}, nil, err
	}

	serviceRequest := &pbd.DumpIdentityListRequest{}

	domainClient := pbd.NewDomainClient(domainConn)
	reply, err := domainClient.DumpIdentityList(ctx, serviceRequest)
	if err != nil {
		return identity.Identity{}, nil, err
	}

	ident, err := identity.ConvertPBItoI(reply.GetIdentity())
	if err != nil {
		return identity.Identity{}, nil, err
	}

	idents, err := identity.ConvertPBItoIMultiple(reply.GetIdentityList())
	if err != nil {
		return identity.Identity{}, nil, err
	}

	return ident, idents, nil
}

func GatherStandardArgs() (port int, domainPort int, logger *log.Logger, err error) {
	portString := os.Args[1]
	port64, err := strconv.ParseInt(portString, 0, 32)
	if err != nil {
		return 0, 0, nil, err
	}
	port = int(port64)

	domainPortString := os.Args[2]
	domainPort64, err := strconv.ParseInt(domainPortString, 0, 32)
	if err != nil {
		return 0, 0, nil, err
	}
	domainPort = int(domainPort64)

	logFileName := os.Args[3]
	if logFileName != "" {
		logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return 0, 0, nil, err
		}
		logger = log.New(logFile, "", log.LstdFlags)
	} else {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	return
}
