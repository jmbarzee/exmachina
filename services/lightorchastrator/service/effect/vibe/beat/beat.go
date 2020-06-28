package beat

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/shared"
)

type Beat struct {
	shared.TimeSpan
}

func (b Beat) Stabalize() Beat {

}

func (b Beat) Materialize() effect.Vibe {

}
