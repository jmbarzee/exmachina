package domain

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

	config := domainConfigExtra{}
	err := toml.Unmarshal(bytes, &config)
	if err != nil {
		return DomainConfig{}, err
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
