package npsub

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/system"
)

func (s *NPSub) SubscribeLights(ctx context.Context) {
	routineName := "SubscribeLights"
	system.LogRoutinef(routineName, "Starting routine")
	ticker := time.NewTicker(time.Second)

Loop:
	for {
		select {
		case <-ticker.C:
			idents, err := s.Service.RPCGetServices(ctx, "lightOrchestrator")
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

			err = s.rpcSubscribeLights(ctx, idents[0])
			if err != nil {
				system.Errorf("Failed call to rpcSubscribeLights: %w", err)
			}
		case <-ctx.Done():
			break Loop
		}
	}
	system.LogRoutinef(routineName, "Stopping routine")
}
