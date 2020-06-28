package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmbarzee/domain/server/identity"
	"github.com/jmbarzee/domain/services/envorchastrator"
	"github.com/jmbarzee/domain/services/exporchastrator"
	"github.com/jmbarzee/domain/services/lightorchastrator"
	"github.com/jmbarzee/domain/services/musicinforet"
	"github.com/jmbarzee/domain/services/nplight"
	"github.com/jmbarzee/domain/services/soundorchastrator"
	"github.com/jmbarzee/domain/services/speaker"
	"github.com/jmbarzee/domain/services/webserver"
)

type (
	Service struct {
		ServiceIdentity identity.ServiceIdentity
		ServiceConfig   ServiceConfig
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
		case "webServer":
			err = webserver.Start(port, d.config.Port)

		case "musicInformationRetrival":
			err = musicinforet.Start(port, d.config.Port)

		case "experienceOrchastrator":
			err = exporchastrator.Start(port, d.config.Port)

		case "lightOrchastrator":
			err = lightorchastrator.Start(port, d.config.Port)
		case "neoPixelLight":
			err = nplight.Start(port, d.config.Port)

		case "soundOrchastrator":
			err = soundorchastrator.Start(port, d.config.Port)
		case "dmlSpeaker":
			err = speaker.Start(port, d.config.Port)

		case "enviornmentOrchastrator":
			err = envorchastrator.Start(port, d.config.Port)
		case "thermostat":
			// err = neopixelbar.Start(port, d.Log)
		case "shade":
			// err = neopixelbar.Start(port, d.Log)
		case "desk":
			// err = neopixelbar.Start(port, d.Log)

		default:
			err = fmt.Errorf("Unknown service! (%s)", config.Name)
		}
		if err != nil {
			goto Unlock
		}
		d.services[config.Name] = Service{
			ServiceIdentity: identity.ServiceIdentity{
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
