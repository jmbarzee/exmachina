package ifaces

// Bender is used by Shifters to change small things over time
type Bender interface {
	Stabalizable

	// Bend returns a value based on f
	Bend(f float64) float64
}
