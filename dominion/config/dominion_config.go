package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jmbarzee/dominion/system"
)

type (
	// DominionConfig holds all the information required to start a Dominion
	DominionConfig struct {
		// Port is the port which the domain will be responding on
		Port int

		// DialTimeout is how long a domain will wait for a grpc.ClientConn to establish
		DialTimeout time.Duration
		// DomainCheck is the length of time after which a dominion will send a heartbeat
		DomainCheck time.Duration
		// ServiceCheck is the length of time after which a dominion check service dependency
		ServiceCheck time.Duration

		// Services is a map of service type to service config
		Services map[string]ServiceDefinition
	}
)

var dominionConfig *DominionConfig

func GetDominionConfig() DominionConfig {
	if dominionConfig == nil {
		system.Panic(errors.New("config.Setup has not ben called"))
	}
	return *dominionConfig
}

func setupDominionConfigFromTOML(configFilePath string) error {
	if dominionConfig != nil {
		return errors.New("config.setupDominionConfigFromTOML has already been called")
	}

	configFile, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}

	config := &DominionConfig{}
	if err = toml.Unmarshal(bytes, config); err != nil {
		return err
	}

	if err = config.check(); err != nil {
		return err
	}

	dominionConfig = config
	return nil
}

func (c DominionConfig) check() error {
	if c.Port == 0 {
		return fmt.Errorf("Port was not set")
	}

	if c.DialTimeout == 0 {
		return fmt.Errorf("ConnectionConfig.DialTimeout was not set")
	}
	if c.DomainCheck == 0 {
		return fmt.Errorf("ConnectionConfig.HeartbeatCheckv was not set")
	}
	if c.ServiceCheck == 0 {
		return fmt.Errorf("ServiceHierarchyConfig.DependencyCheck was not set")
	}
	return nil
}

func (c DominionConfig) String() string {
	dumpMsg := "\tPort: " + strconv.Itoa(c.Port) + "\n" +
		"\tDialTimeout: " + c.DialTimeout.String() + "\n" +
		"\tDomainCheck: " + c.DomainCheck.String() + "\n" +
		"\tServiceCheck: " + c.ServiceCheck.String() + "\n" +
		"\tServices: {\n"
	for serviceType, serviceConfig := range c.Services {
		dumpMsg += "\t\t" + serviceType + ": " + serviceConfig.String() + ",\n"
	}
	dumpMsg += "\t}\n"

	return dumpMsg
}
