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
	DomainConfig struct {
		// UUID is a unique identifier for a domain
		UUID string
		// Port is the port which the domain will be responding on
		Port int

		// Traits is the traits possesed by the domain.
		Traits []string

		// DialTimeout is how long a domain will wait for a grpc.ClientConn to establish
		DialTimeout time.Duration
		// IsolationCheck is the duration between isolation checks
		IsolationCheck time.Duration
		// ServiceCheck is the length of time after which a domain will send a hearbeat
		ServiceCheck time.Duration
	}
)

var domainConfig *DomainConfig

func GetDomainConfig() DomainConfig {
	if domainConfig == nil {
		system.Panic(errors.New("config.Setup has not ben called"))
	}
	return *domainConfig
}

func SetupFromTOML(configFilePath string) error {
	if domainConfig != nil {
		return errors.New("config.setupFromTOML has already been called")
	}

	configFile, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}

	config := &DomainConfig{}
	if err = toml.Unmarshal(bytes, config); err != nil {
		return err
	}

	if err = config.check(); err != nil {
		return err
	}

	domainConfig = config
	return nil
}

func (c DomainConfig) check() error {

	if c.UUID == "" {
		return fmt.Errorf("UUID was not set")
	}
	if len(c.Traits) == 0 {
		return fmt.Errorf("Traits were not set")
	}
	if c.Port == 0 {
		return fmt.Errorf("Port was not set")
	}

	if c.DialTimeout == 0 {
		return fmt.Errorf("DialTimeout was not set")
	}
	if c.IsolationCheck == 0 {
		return fmt.Errorf("IsolationCheck was not set")
	}
	if c.ServiceCheck == 0 {
		return fmt.Errorf("ServiceCheck was not set")
	}

	return nil
}

func (c DomainConfig) String() string {

	dumpMsg := "\tUUID: " + c.UUID + "\n"

	dumpMsg += "\tTraits: ["
	for _, trait := range c.Traits {
		dumpMsg += trait + ", "
	}
	dumpMsg += "]\n"

	dumpMsg += "\tPort: " + strconv.Itoa(c.Port) + "\n" +
		"\tDialTimeout: \n" + strconv.Itoa(int(c.DialTimeout)) +
		"\tServiceCheck: \n" + strconv.Itoa(int(c.ServiceCheck)) +
		"\tIsolationCheck: \n" + strconv.Itoa(int(c.IsolationCheck))
	return dumpMsg
}
