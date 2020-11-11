package ifaces

// Vibe is a heavy abstraction correlating to general feelings in music
type Vibe interface {
	Palette

	Stabalizable

	// Duplicate creates a copy of a vibe and insures that
	// the dupliacted vibe will stabalize/materialize differently
	Duplicate() Vibe

	// Stabilize locks in part of the visual representation of a vibe.
	Stabilize() Vibe

	// Materialize locks all remaining unlocked visuals of a vibe
	// then returns the resulting effects
	Materialize() []Effect
}
