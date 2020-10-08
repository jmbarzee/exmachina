package shared

// StabilizeFunc freezes some trait
type StabilizeFunc func(p Palette) error

type stabilizable interface {
	// GetStabilizeFuncs returns StabilizeFunc for all remaining unstablaized traits
	GetStabilizeFuncs() []StabilizeFunc
}
