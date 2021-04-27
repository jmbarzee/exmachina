package span

import "time"

// Span is represents anything that starts and Ends
type Span struct {
	StartTime time.Time
	EndTime   time.Time
}

// Start returns the Start time
func (s Span) Start() time.Time { return s.StartTime }

// End returns the End time
func (s Span) End() time.Time { return s.EndTime }
