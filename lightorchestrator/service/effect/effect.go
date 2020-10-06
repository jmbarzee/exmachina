package effect

import (
	"time"

	"github.com/jmbarzee/dominion/services/lightorchestrator/service/shared"
)

// Effect is a light abstraction representing paterns of colors
type Effect interface {
	// Render will produce a slice of lights based on the time and properties of lights
	Render(t time.Time, lights []shared.Light) []shared.Light
	// Priority solves rendering issues
	Priotity() int

	// Start returns the Start time
	Start() time.Time
	// End returns the End time
	End() time.Time
}
