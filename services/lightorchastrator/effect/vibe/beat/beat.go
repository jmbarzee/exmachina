package beat

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/effect/shared"
)

type Beat struct {
	shared.TimeSpan
}

func (b Beat) Stabalize() Beat {

}

func (b Beat) Materialize() effect.Vibe {

}
