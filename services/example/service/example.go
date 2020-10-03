package service

import (
	"context"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	pb "github.com/jmbarzee/dominion/services/example/grpc"
	"github.com/jmbarzee/dominion/system"
)

type ExampleService struct {
	*service.Service
}

func NewExampleService(config config.ServiceConfig) (ExampleService, error) {
	service, err := service.NewService(config)
	if err != nil {
		return ExampleService{}, err
	}

	example := ExampleService{
		Service: service,
	}

	pb.RegisterExampleServiceServer(service.Server, example)
	return example, nil
}

func (s *ExampleService) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(s.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go s.exampleRoutine(ctx)

	return s.Service.HostService(ctx)
}
