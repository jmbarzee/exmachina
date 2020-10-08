package shared

import (
	"time"
)

// Shifter is used by Effects and Painters to provide changing values based on time
type Shifter interface {
	// Shift returns a value representing some change or shift
	Shift(t time.Time) float32
	stabilizable
}
