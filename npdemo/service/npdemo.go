package service

import (
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/lightorchestrator/clients/nplight"
)

type (
	NPDemo struct {
		*nplight.NPLight
	}
)

const (
	size = 120 //TODO @jmbarzee fetch from config
)

func NewNPDemo(config config.ServiceConfig) (NPDemo, error) {
	sub, err := nplight.NewNPLight(config, size)
	if err != nil {
		return NPDemo{}, err
	}
	return NPDemo{sub}, nil
}
