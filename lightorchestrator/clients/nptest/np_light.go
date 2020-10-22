package nptest

import (
	"context"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/clients/npsub"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second // displayFPS

	pixelsPerLight = 10
)

type (
	NPTest struct {
		*npsub.NPSub
		Window *pixelgl.Window
	}
)

func NewNPTest(config config.ServiceConfig, size int) (*NPTest, error) {
	sub, err := npsub.NewNPSub(config, size)
	if err != nil {
		return nil, err
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, float64(pixelsPerLight*sub.Size), pixelsPerLight*2),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return &NPTest{
		NPSub:  sub,
		Window: win,
	}, nil
}

func (l *NPTest) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(l.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go l.SubscribeLights(ctx)
	go system.RoutineOperation(ctx, "UpdateLights", displayRate, l.updateLights)

	return l.Service.HostService(ctx)
}
