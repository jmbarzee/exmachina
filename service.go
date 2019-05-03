package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmbarzee/domain/services/light"
	"github.com/jmbarzee/domain/services/lightfeed"
	"github.com/jmbarzee/domain/services/musicinfo"
)

type (
	Service struct {
		ServiceIdentity ServiceIdentity
		ServiceConfig   ServiceConfig
	}

	ServiceIdentity struct {
		Port        int
		LastContact time.Time
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

func (s *Service) Status() error {
	return errors.New("Unimplemented!")
}

func (s ServiceConfig) String() string {
	msg := "(" + s.Name + ", " +
		string(s.Priority) + ", "

	msg += "["
	for _, dependency := range s.Depends {
		msg += dependency + ", "
	}
	msg += "]"

	msg += ", "

	msg += "["
	for _, trait := range s.Traits {
		msg += trait + ", "
	}
	msg += "]"

	msg += ")"
	return msg
}

func (d *Domain) serviceConfigFromName(serviceName string) (ServiceConfig, error) {
	for _, serviceConfig := range d.config.Services {
		if serviceConfig.Name != serviceName {
			continue
		} else {
			return serviceConfig, nil
		}
	}
	return ServiceConfig{}, fmt.Errorf("serviceConfig '%v' not found", serviceName)
}

func (d *Domain) hasTraitsForService(serviceConfig ServiceConfig) bool {
	hasAllTraits := true
	for _, trait := range serviceConfig.Traits {
		hasAllTraits = hasAllTraits && d.hasTrait(trait)
	}
	return hasAllTraits
}

func (d *Domain) getProficiencyForService(serviceConfig ServiceConfig) int32 {
	// TODO @jmbarzee make this more intelligent
	return 1
}

func (d *Domain) startService(config ServiceConfig) error {
	if !d.hasTraitsForService(config) {
		return errors.New("tried to start service without needed traits")
	}

	var err error
	d.debugf(debugLocks, "startService() pre-lock(%v)\n", "servicesLock")
	d.servicesLock.Lock()
	{
		d.debugf(debugLocks, "startService() in-lock(%v)\n", "servicesLock")

		port := d.config.Port + len(d.services) // TODO @jmbarzee Hacky, find a way to make more durable
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
		d.services[config.Name] = Service{
			ServiceIdentity: ServiceIdentity{
				Port:        port,
				LastContact: time.Now(),
			},
			ServiceConfig: config}
		d.Logf("Started service: \"%v\" at port:%v", config.Name, port)

	Unlock:
	}
	d.servicesLock.Unlock()
	d.debugf(debugLocks, "startService() post-lock(%v)\n", "servicesLock")
	return err
}
