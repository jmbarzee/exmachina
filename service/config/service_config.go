package config

import (
	"net"
	"os"
	"strconv"
)

type ServiceConfig struct {
	DominionIP   net.IP
	DominionPort int
	ServicePort  int
	ServiceType  string
}

func FromEnv(serviceType string) (config ServiceConfig, err error) {
	dominionIPString := os.Args[1]
	dominionIP := net.ParseIP(dominionIPString)

	dominionPortString := os.Args[2]
	dominionPort64, err := strconv.ParseInt(dominionPortString, 0, 32)
	if err != nil {
		return ServiceConfig{}, err
	}
	dominionPort := int(dominionPort64)

	servicePortString := os.Args[3]
	servicePort64, err := strconv.ParseInt(servicePortString, 0, 32)
	if err != nil {
		return ServiceConfig{}, err
	}
	servicePort := int(servicePort64)

	return ServiceConfig{
		DominionIP:   dominionIP,
		DominionPort: dominionPort,
		ServicePort:  servicePort,
		ServiceType:  serviceType,
	}, nil
}
