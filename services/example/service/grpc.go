package service

import (
	"context"
	"errors"

	pb "github.com/jmbarzee/dominion/services/example/grpc"
)

func (s ExampleService) ExampleRPC(ctx context.Context, request *pb.ExampleRPCRequest) (*pb.ExampleRPCReply, error) {
	return nil, errors.New("Unimplemented")
}
