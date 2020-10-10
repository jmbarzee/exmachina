package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
)

// Effect is a light abstraction representing paterns of colors
type Effect interface {
	// Render will produce a slice of lights based on the time and properties of lights
	Render(t time.Time, lights []light.Light) []light.Light
	// Priority solves rendering issues
	Priotity() int

	// Start returns the Start time
	Start() time.Time
	// End returns the End time
	End() time.Time
}
