package service

import (
	"fmt"
	"path"

	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
)

type (
	Service struct {
		identity.ServiceIdentity

		DominionIdentity identity.DominionIdentity
	}
)

func NewService(config config.ServiceConfig) (*Service, error) {
	logFilePath := path.Join("/usr/local/dominion/logs", config.ServiceType+".log")
	if err := system.Setup(logFilePath); err != nil {
		return nil, err
	}

	// Initialize IP
	ip, err := system.GetOutboundIP()
	if err != nil {
		return nil, fmt.Errorf("failed to find Local IP: %v\n", err.Error())
	}

	return &Service{
		ServiceIdentity: identity.ServiceIdentity{
			Type: config.ServiceType,
			Address: identity.Address{
				IP:   ip,
				Port: config.ServicePort,
			},
		},
		DominionIdentity: identity.DominionIdentity{
			Address: identity.Address{
				IP:   config.DominionIP,
				Port: config.DominionPort,
			},
		},
	}, nil
}
