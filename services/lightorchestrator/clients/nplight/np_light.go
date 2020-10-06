package nplight

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/services/lightorchestrator/clients/nplight/lightplan"
	"github.com/jmbarzee/dominion/system"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second // displayFPS

	gpioPin    = 18
	brightness = 255
)

type (
	NPLight struct {
		*service.Service

		Size      int
		LightPlan lightplan.LightPlan
	}
)

func NewNPLight(config config.ServiceConfig, size int) (*NPLight, error) {
	service, err := service.NewService(config)
	if err != nil {
		return nil, err
	}

	return &NPLight{
		Service:   service,
		Size:      size,
		LightPlan: lightplan.NewLightPlan(),
	}, nil
}

func (l *NPLight) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(l.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go l.subscribeLights(ctx)
	go l.displayLights(ctx)

	return l.Service.HostService(ctx)
}
