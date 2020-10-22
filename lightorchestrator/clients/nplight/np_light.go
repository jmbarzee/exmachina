package nplight

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/clients/npsub"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second / displayFPS

	gpioPin    = 18
	brightness = 255
)

type (
	NPLight struct {
		*npsub.NPSub
		Strip *ws2811.WS2811
	}
)

func NewNPLight(config config.ServiceConfig, size int) (*NPLight, error) {
	sub, err := npsub.NewNPSub(config, size)
	if err != nil {
		return nil, err
	}

	opt := ws2811.DefaultOptions
	opt.Channels[0].StripeType = ws2811.SK6812StripRGBW
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = sub.Size
	opt.Channels[0].GpioPin = gpioPin

	strip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		system.Panic(err)
	}

	err = strip.Init()
	if err != nil {
		system.Panic(err)
	}

	return &NPLight{
		NPSub: sub,
		Strip: strip,
	}, nil
}

func (l *NPLight) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(l.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go l.SubscribeLights(ctx)
	go system.RoutineOperation(ctx, "UpdateLights", displayRate, l.updateLights)
	defer l.Strip.Fini()

	return l.Service.HostService(ctx)
}
