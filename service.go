package domain

import (
	"errors"
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
		d.Logf("Started service: \"%v\" at port:%v", config.Name, port)

	Unlock:
	}
	d.ServicesLock.Unlock()
	d.debugf(debugLocks, "startService() post-lock(%v)\n", "ServicesLock")
	return err
}
