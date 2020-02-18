package server

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// ConfigureFromYAML produces a default configuration which can be passed to New()
func ConfigFromYAML(bytes []byte) (DomainConfig, error) {
	return DomainConfig{}, errors.New("Unimplemented!")
}

// ConfigFromTOML produces a default configuration which can be passed to New()
func ConfigFromTOML(bytes []byte) (DomainConfig, error) {
	type (
		domainConfigExtra struct {
			LogFileName string
			DomainConfig
		}
	)

	// Comprehend TOML
	config := domainConfigExtra{}
	err := toml.Unmarshal(bytes, &config)
	if err != nil {
		return DomainConfig{}, err
	}

	// Duplicate Service names
	for serviceName, serviceConfig := range config.Services {
		// Strange copy magic...
		// can't modify the struct member directly, so modify the copy from the range, then set.
		serviceConfig.Name = serviceName
		config.Services[serviceName] = serviceConfig
	}

	// Initialize LogFile
	var logger *log.Logger
	if config.LogFileName != "" {
		logFile, err := os.OpenFile(config.LogFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return DomainConfig{}, err
		}
		logger = log.New(logFile, "", log.LstdFlags)
	} else {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	config.Log = logger

	return config.DomainConfig, nil
}
