package vibe

import (
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/color"
	effect "github.com/jmbarzee/dominion/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/shared"
)

type BasicVibe struct {
	shared.TimeSpan
}

func (b BasicVibe) Stabalize() Vibe {
	return b
}

func (b BasicVibe) Materialize() []effect.Effect {
	return []effect.Effect{
		effect.Solid{
			BasicEffect: effect.BasicEffect{
				TimeSpan: b.TimeSpan,
				Rank:     0,
			},
			Color: effect.NewColorFuncStatic(color.WhiteNatural),
		},
	}
}
