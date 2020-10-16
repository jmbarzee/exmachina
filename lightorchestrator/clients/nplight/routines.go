package nplight

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/system"
	ws2811 "github.com/jmbarzee/rpi_ws281x/golang/ws2811"
)

func (l *NPLight) displayLights(ctx context.Context) {
	routineName := "DisplayLights"
	system.LogRoutinef(routineName, "Starting routine")
	defer ws2811.Fini()
	err := ws2811.Init(gpioPin, l.Size, brightness)
	if err != nil {
		system.Panic(err)
	}

	ticker := time.NewTicker(displayRate)

Loop:
	for {
		select {
		case t := <-ticker.C:
			// Advance the light plan
			next := l.LightPlan.Advance(t)
			if next != nil {
				for i, wrgb := range next.Lights {
					ws2811.SetLed(i, wrgb)
				}
				ws2811.Render()
			}

		case <-ctx.Done():
			break Loop
		}
	}
	system.LogRoutinef(routineName, "Stopping routine")
}

func (l *NPLight) subscribeLights(ctx context.Context) {
	routineName := "SubscribeLights"
	system.LogRoutinef(routineName, "Starting routine")
	ticker := time.NewTicker(time.Second)

Loop:
	for {
		select {
		case <-ticker.C:
			idents, err := l.Service.RPCGetServices(ctx, "lightOrchestrator")
			if err != nil {
				system.Logf("Error locating lightOrchestrator: %v", err.Error())
				continue
			}
			if len(idents) > 1 {
				system.Logf("Found multiple lightOrchestrator, %v", idents)
				continue
			}
			if len(idents) < 1 {
				continue
			}
			system.LogRoutinef(routineName, "Found new orchestrator")

			err = l.rpcSubscribeLights(ctx, idents[0])
			if err != nil {
				system.Errorf("Failed call to rpcSubscribeLights: %w", err)
			}
		case <-ctx.Done():
			break Loop
		}
	}
	system.LogRoutinef(routineName, "Stopping routine")
}
