package example

import (
	"context"

	"github.com/jmbarzee/dominion/system"
)

func (s ExampleService) exampleRoutine(ctx context.Context) {
	routineName := "exampleRoutine"
	system.LogRoutinef(routineName, "Starting routine")
	select {
	case <-ctx.Done():
	}
	system.LogRoutinef(routineName, "Stopping routine")
}
