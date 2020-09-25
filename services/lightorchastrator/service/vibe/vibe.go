package vibe

import (
	"time"

	"github.com/jmbarzee/dominion/services/lightorchastrator/service/effect"
)

// Vibe is a heavy abstraction correlating to general feelings in music
type Vibe interface {
	// Stabalize locks in part of the visual representation of a vibe.
	Stabalize() Vibe
	// Materialize locks all remaining unlocked visuals of a vibe
	// then returns the resulting effects
	Materialize() []effect.Effect

	// Start returns the Start time
	Start() time.Time
	// End returns the End time
	End() time.Time
}
