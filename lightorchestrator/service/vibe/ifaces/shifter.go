package ifaces

import "time"

// Shifter is used by Painters to change small things over time
type Shifter interface {
	Stabalizable

	// Shift returns a value representing some change or shift based on t
	Shift(t time.Time) float64
}
