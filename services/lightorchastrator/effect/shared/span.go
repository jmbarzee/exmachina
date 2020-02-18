package shared

import "time"

// TimeSpan is represents anything that starts and Ends
type TimeSpan struct {
	StartTime time.Time
	EndTime   time.Time
}

// Start returns the Start time
func (s TimeSpan) Start() time.Time { return s.StartTime }

// End returns the End time
func (s TimeSpan) End() time.Time { return s.EndTime }
