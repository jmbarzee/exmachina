package ifaces

type Stabalizable interface {
	// GetStabilizeFuncs returns a function for all remaining unstablaized traits
	GetStabilizeFuncs() []func(p Palette)
}
