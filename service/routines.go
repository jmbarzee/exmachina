package service

import (
	"context"
	"fmt"
	"net"

	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Service) HostService(ctx context.Context) error {
	routineName := "HostService"
	system.LogRoutinef(routineName, "Starting routine")

	address := fmt.Sprintf("%s:%v", "", s.Address.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("hostService() Failed to listen: %w", err)
	}

	server := grpc.NewServer()
	pb.RegisterServiceServer(server, s)
	// Register reflection service on gRPC server.
	go func() {
		<-ctx.Done()
		server.GracefulStop()
		system.LogRoutinef(routineName, "Stopped grpc server gracefully.")
	}()

	reflection.Register(server)
	err = server.Serve(lis)
	system.LogRoutinef(routineName, "Stopping routine")
	return err
}
