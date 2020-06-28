package service

import (
	"log"

	"github.com/jmbarzee/domain/services/lightorchastrator/clients/nplight"
)

type (
	NPLight struct {
		*nplight.Subscriber
	}
)

const (
	size = 30 //TODO @jmbarzee fetch from config
)

func NewNPLight(port int, domainPort int, logger *log.Logger) NPLight {
	return NPLight{
		Subscriber: nplight.NewSubscriber(port, domainPort, logger, size),
	}
}
