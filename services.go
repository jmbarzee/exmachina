package domain

import (
	"fmt"

	"github.com/jmbarzee/domain/services/light"
	"github.com/jmbarzee/domain/services/lightfeed"
	"github.com/jmbarzee/domain/services/musicinfo"
)

func (d *Domain) startRequiredServices() {
	for _, serviceConfig := range d.config.Services {
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
		d.startService(serviceConfig)
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

func (d *Domain) startService(config ServiceConfig) error {
	var err error
	d.debugf(debugLocks, "startService() pre-lock(%v)\n", "ServicesLock")
	d.ServicesLock.Lock()
	{
		d.debugf(debugLocks, "startService() in-lock(%v)\n", "ServicesLock")

		port := d.config.Port + len(d.services)
		if _, ok := d.services[config.Name]; ok {
			err = fmt.Errorf("Service already exists! (%s)", config.Name)
			goto Unlock
		}

		switch config.Name {
		case "light":
			err = light.Start(port, d.Log)
		case "lightFeed":
			err = lightfeed.Start(port, d.Log)
		case "musicInfo":
			err = musicinfo.Start(port, d.Log)
		default:
			err = fmt.Errorf("Unknown service! (%s)", config.Name)
		}
		if err != nil {
			goto Unlock
		}

		d.services[config.Name] = Service{Port: port, ServiceConfig: config}

	Unlock:
	}
	d.ServicesLock.Unlock()
	d.debugf(debugLocks, "startService() post-lock(%v)\n", "ServicesLock")
	return err
}
