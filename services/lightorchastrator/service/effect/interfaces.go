package effect

import (
	"time"
)

type (
	Allocater interface {
		Allocate(Vibe)
	}

	Device interface {
		Allocater
		Render(time.Time) []*Light
		PruneEffects(time.Time)
	}

	Vibe interface {
		Start() time.Time
		End() time.Time
		Stabalize() Vibe
		Materialize() []Effect
	}

	Effect interface {
		Start() time.Time
		End() time.Time
		Priotity() int
		Render(time.Time, []*Light)
	}
)
