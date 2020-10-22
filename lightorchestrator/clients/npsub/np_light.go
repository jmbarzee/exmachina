package npsub

import (
	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/lightorchestrator/clients/npsub/lightplan"
)

type (
	NPSub struct {
		*service.Service

		Size      int
		LightPlan lightplan.LightPlan
	}
)

func NewNPSub(config config.ServiceConfig, size int) (*NPSub, error) {
	service, err := service.NewService(config)
	if err != nil {
		return nil, err
	}

	return &NPSub{
		Service:   service,
		Size:      size,
		LightPlan: lightplan.NewLightPlan(),
	}, nil
}
