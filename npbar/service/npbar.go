package service

import (
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/lightorchestrator/clients/nplight"
)

type (
	NPBar struct {
		*nplight.NPLight
	}
)

const (
	size = 30 //TODO @jmbarzee fetch from config
)

func NewNPBar(config config.ServiceConfig) (NPBar, error) {
	sub, err := nplight.NewNPLight(config, size)
	if err != nil {
		return NPBar{}, err
	}
	return NPBar{sub}, nil
}
