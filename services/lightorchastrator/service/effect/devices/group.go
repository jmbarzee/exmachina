package devices

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/shared"
)

type GroupOption struct {
	Groups []Group
}

func (o GroupOption) Allocate(feeling effect.Vibe) {
	groupNum := shared.RepeatableOption(feeling.Start(), len(o.Groups))
	o.Groups[groupNum].Allocate(feeling)
}

type Group struct {
	Devices []effect.Allocater
}

func (g Group) Allocate(feeling effect.Vibe) {
	newFeeling := feeling.Stabalize()

	for _, device := range g.Devices {
		device.Allocate(newFeeling)
	}
}
