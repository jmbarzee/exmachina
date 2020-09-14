package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jmbarzee/dominion/system"
)

type (
	ServicesConfig struct {
		Services map[string]ServiceDefinition
	}

	// ServiceDefinition defines a single service in the service hiarchy
	ServiceDefinition struct {
		// Priority is the priority of the service
		Priority Priority
		// Dependencies is the list of service types which this service depends on
		Dependencies []string
		// Traits is the list of triats required by a domain to be able to run a service
		Traits []string
	}

	Priority string
)

var servicesConfig *ServicesConfig

func GetServicesConfig() ServicesConfig {
	if servicesConfig == nil {
		system.Panic(errors.New("config.Setup has not ben called"))
	}
	return *servicesConfig
}

func setupServicesConfigFromTOML(configFilePath string) error {
	if servicesConfig != nil {
		return errors.New("config.setupServicesConfigFromTOML has already been called")
	}

	configFile, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}

	config := &ServicesConfig{}
	if err = toml.Unmarshal(bytes, config); err != nil {
		return err
	}

	servicesConfig = config
	return nil
}

func (c ServicesConfig) GetRequiredServices() map[string]ServiceDefinition {
	requiredServices := make(map[string]ServiceDefinition, 0)
	for serviceType, serviceDef := range c.Services {
		if serviceDef.IsRequired() {
			requiredServices[serviceType] = serviceDef
		}
	}
	return requiredServices
}

const (
	Required   Priority = "required"
	Dependency Priority = "dependency"
)

func (s ServiceDefinition) String() string {
	msg := "(" + string(s.Priority) + ", ["

	for _, dependency := range s.Dependencies {
		msg += dependency + ", "
	}
	msg += "], ["

	for _, trait := range s.Traits {
		msg += trait + ", "
	}
	msg += "])"

	return msg
}

func (s ServiceDefinition) IsRequired() bool {
	return s.Priority == Required
}
