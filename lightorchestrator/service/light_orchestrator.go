package service

import (
	"context"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	pb "github.com/jmbarzee/dominion/services/lightorchestrator/grpc"
	"github.com/jmbarzee/dominion/system"
)

type LightOrch struct {
	*service.Service

	Subscribers *SubscriberList

	DeviceHierarchy *DeviceNodeTree
}

func NewLightOrch(config config.ServiceConfig) (*LightOrch, error) {
	service, err := service.NewService(config)
	if err != nil {
		return nil, err
	}

	subscriberList, deviceNodeTree := NewStructs()

	lightOrch := &LightOrch{
		Service:         service,
		Subscribers:     subscriberList,
		DeviceHierarchy: deviceNodeTree,
	}

	pb.RegisterLightOrchestratorServer(service.Server, lightOrch)
	return lightOrch, nil
}

func (l *LightOrch) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(l.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go l.orchastrate(ctx)

	return l.Service.HostService(ctx)
}
