package config

import (
	"net"
	"os"
	"strconv"
	"time"
)

type ServiceConfig struct {
	DominionIP   net.IP
	DominionPort int
	DomainUUID   string
	ServicePort  int
	ServiceType  string
}

const DefaultServiceDialTimeout = time.Millisecond * 100

func FromEnv(serviceType string) (config ServiceConfig, err error) {
	dominionIPString := os.Args[1]
	dominionIP := net.ParseIP(dominionIPString)

	dominionPortString := os.Args[2]
	dominionPort64, err := strconv.ParseInt(dominionPortString, 0, 32)
	if err != nil {
		return ServiceConfig{}, err
	}
	dominionPort := int(dominionPort64)

	domainUUID := os.Args[3]

	servicePortString := os.Args[4]
	servicePort64, err := strconv.ParseInt(servicePortString, 0, 32)
	if err != nil {
		return ServiceConfig{}, err
	}
	servicePort := int(servicePort64)

	return ServiceConfig{
		DominionIP:   dominionIP,
		DominionPort: dominionPort,
		DomainUUID:   domainUUID,
		ServicePort:  servicePort,
		ServiceType:  serviceType,
	}, nil
}
