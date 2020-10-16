package service

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
)

const (
	displayFPS                = 1
	displayRate time.Duration = time.Second / displayFPS
)

func (l *LightOrch) orchastrate(ctx context.Context) {
	routineName := "orchastrate"
	system.LogRoutinef(routineName, "Starting routine")
	ticker := time.NewTicker(displayRate)

Loop:
	for {

		select {
		case <-ticker.C:
			l.Subscribers.Range(func(sub Subscriber) bool {
				if err := sub.DispatchRender(time.Now()); err != nil {
					system.Errorf("Failed to dispatch Render: %w", err)
				}
				return true
			})
		case <-ctx.Done():
			break Loop
		}
	}

	system.LogRoutinef(routineName, "Stopping routine")
}

func (l *LightOrch) subscribeVibes(ctx context.Context) {
	routineName := "orchastrate"
	system.LogRoutinef(routineName, "Starting routine")

	tickLength := time.Second * 20
	ticker := time.NewTicker(tickLength)

Loop:
	for {

		select {
		case t := <-ticker.C:
			v := &vibe.Basic{
				Span: span.Span{
					StartTime: t.Add(tickLength),
					EndTime:   t.Add(tickLength * 2),
				},
			}
			l.DeviceHierarchy.Allocate(v)
		case <-ctx.Done():
			break Loop
		}
	}
	system.LogRoutinef(routineName, "Stopping routine")
}
