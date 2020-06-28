package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/blang/semver"
	"github.com/jmbarzee/domain/server/system"
)

const Version = "0.1.0"

type (
	Domain struct {
		// System contains the halt and log of the domain.
		// System handles Signals with `go system.HandleSignals`.
		// System offers `system.Logf` and `system.Panic`.
		system.System

		// services is the list of services the domain currently offers as
		//    serviceName -> Service
		services     map[string]Service
		servicesLock sync.Mutex

		// electionMap stores the current running elections as
		//     serviceName -> Election
		elections     map[string]*Election
		electionsLock sync.Mutex

		// peerMap stores the members of a Dominion in a wrapped sync.map as
		//     GUID -> peer
		peerMap *peerMap

		// config does what it is. See DomainConfig for a clear understanding.
		// It is unchanged after NewDomain() completes
		config DomainConfig
	}
)

// NewDomain creates a new Domain, to correctly build the Domain, just initilize
func NewDomain(config DomainConfig) (*Domain, error) {

	// Check config
	err := config.Check()
	if err != nil {
		return nil, err
	}

	// Initialize Version
	config.Version, err = semver.Parse(Version)
	if err != nil {
		return nil, fmt.Errorf("failed to semver.Parse(%v): %v\n", Version, err.Error())
	}

	// Initialize IP
	config.IP, err = getOutboundIP()
	if err != nil {
		return nil, fmt.Errorf("failed to find Local IP: %v\n", err.Error())
	}

	return &Domain{
		config:    config,
		services:  make(map[string]Service, 0),
		elections: make(map[string]*Election, 0),
		peerMap: &peerMap{
			sMap: &sync.Map{},
		},
	}, nil
}

func (d Domain) Run(ctx context.Context) {
	// Initilize System (logging, context, and signals)
	var systemCtx context.Context
	systemCtx, d.System = system.NewSystem(ctx, d.config.Log)

	// Dump Stats
	startMsg := "I seek to join the Dominion\n" +
		d.config.Dump() +
		"The Dominion ever expands!\n"

	d.Logf(startMsg)

	// Start Auto Connecting Routines
	go d.watchIsolation(systemCtx)
	go d.listenForBroadcasts(systemCtx)
	go d.buildDomain(systemCtx)

	// Start Services
	d.startRequiredServices()
	d.watchServicesDepnedencies(systemCtx)
}

const (
	debugRoutines       = "Routine"
	debugRPCs           = "RPC"
	debugDomain         = "Domain"
	debugLocks          = "Locks"
	debugLocksElections = "LocksElections"
	debugFatal          = "Fatal"
	debugDefault        = "Default"
)

func (d Domain) debugf(class, fmt string, args ...interface{}) {
	const (
		debug             = true
		logRoutines       = false
		logRPCs           = false
		logDomain         = false
		logLocks          = false
		logLocksElections = false
		logFatal          = true
	)

	if !debug {
		return
	}

	switch class {
	case debugRoutines:
		if logRoutines {
			d.LogfGUID(fmt, args...)
		}
	case debugRPCs:
		if logRPCs {
			d.LogfGUID(fmt, args...)
		}
	case debugDomain:
		if logDomain {
			d.LogfGUID(fmt, args...)
		}
	case debugLocks:
		if logLocks {
			d.LogfGUID(fmt, args...)
		}
	case debugLocksElections:
		if logLocksElections {
			d.LogfGUID(fmt, args...)
		}
	case debugFatal:
		if logFatal {
			d.LogfGUID(fmt, args...)
		}
	case debugDefault:
		d.LogfGUID(fmt, args...)
	default:
		d.LogfGUID("debugf, unknown class:"+fmt, args...)
	}
}

func (d Domain) LogfGUID(fmt string, args ...interface{}) {
	d.Logf(" - "+d.config.UUID+" - "+fmt, args...)
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
