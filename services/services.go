package services

import (
	"fmt"

	"github.com/jmbarzee/domain/services/light"
	"github.com/jmbarzee/domain/services/lightfeed"
	"github.com/jmbarzee/domain/services/musicinfo"
)

type (
	Service struct {
		Port          int
		ServiceConfig ServiceConfig
	}

	ServiceConfig struct {
		Name     string
		Priority Priority
		Depends  []string
		Traits   []string
	}
)

type Priority string

const (
	Required   Priority = "required"
	Dependency Priority = "dependency"
)

func (d *Domain) startRequiredServices() {
	for serviceName, serviceConfig := range d.config.Services {
		if serviceConfig.Priority != Required {
			continue
		}

		// Check to see if domain posses the nessecarry traits
		hasAllTraits := true
		for _, trait := range serviceConfig.Traits {
			hasAllTraits = hasAllTraits && d.hasTrait(trait)
		}
		if !hasAllTraits {
			continue
		}

		// Domain CAN and MUST start the service
		d.startService(serviceName)
	}
}

func (d *Domain) hasTrait(trait string) bool {
	for _, ownTrait := range d.config.Traits {
		if ownTrait == trait {
			return true
		}
	}
	return false
}

func (d *Domain) startService(serviceName string) error {

	var err error
	d.debugf(debugLocks, "startService() pre-lock(%v)\n", "ServicesLock")
	d.ServicesLock.Lock()
	{
		port := d.config.Port + len(d.Services)
		d.debugf(debugLocks, "startService() in-lock(%v)\n", "ServicesLock")
		for name := range d.Services {
			if name == serviceName {
				err = fmt.Errorf("Service already exists! (%s)", serviceName)
				goto Unlock
			}
		}

		switch serviceName {
		case "light":
			light.Start(port)
		case "lightFeed":
			lightfeed.Start(port)
		case "musicInfo":
			musicinfo.Start(port)
		}

	Unlock:
	}
	d.ServicesLock.Unlock()
	d.debugf(debugLocks, "startService() post-lock(%v)\n", "ServicesLock")
	return err
}
