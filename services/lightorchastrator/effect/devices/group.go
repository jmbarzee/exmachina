package devices

import (
	"github.com/jmbarzee/domain/services/lightorchastrator/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/effect/shared"
)

type GroupOption struct {
	Groups []Group
}

func (o GroupOption) Allocate(feeling effect.Feeling) {
	groupNum := shared.RepeatableOption(feeling.Start(), len(o.Groups))
	o.Groups[groupNum].Allocate(feeling)
}

type Group struct {
	Devices []effect.Allocater
}

func (g Group) Allocate(feeling effect.Feeling) {
	newFeeling := feeling.Stabalize()

	for _, device := range g.Devices {
		device.Allocate(newFeeling)
	}
}
