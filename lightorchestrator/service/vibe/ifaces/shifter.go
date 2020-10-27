package ifaces

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
)

// Shifter is used by Painters to change small things over time
type Shifter interface {
	Stabalizable

	// Shift returns a value representing some change or shift based on t and l
	Shift(t time.Time, l light.Light) float64
}
