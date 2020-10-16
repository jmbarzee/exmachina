package ifaces

import "time"

// Span is any object which has a beginning and end in time
type Span interface {
	// Start returns the Start time
	Start() time.Time
	// End returns the End time
	End() time.Time
}
