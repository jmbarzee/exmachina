package domain

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/blang/semver"
	"github.com/jmbarzee/domain/system"
)

const Version = "0.1.0"

type (
	Domain struct {
		// System contains the halt and log of the domain.
		// System handles Signals with `go system.HandleSignals`.
		// System offers `system.Logf` and `system.Panic`.
		system.System

		// services is the list of services the domain currently offers
		services     map[string]Service
		ServicesLock sync.Mutex

		// peerMap stores the members of a Dominion in a wrapped sync.map
		peerMap *peerMap

		// config does what it is. See DomainConfig for a clear understanding.
		// It is unchanged after NewDomain() completes
		config DomainConfig
	}
)

// NewDomain creates a new Domain, to correctly build the Domain, just initilize
func NewDomain(ctx context.Context, config DomainConfig) (*Domain, error) {

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

	d := &Domain{
		config: config,
		peerMap: &peerMap{
			sMap: &sync.Map{},
		},
		services: make(map[string]Service, 0),
	}

	// Initilize System (logging, context, and signals)
	var systemCtx context.Context
	systemCtx, d.System = system.NewSystem(ctx, config.Log)

	// Start Auto Connecting Routines
	go d.watchIsolation(systemCtx)
	go d.listenForBroadcasts(systemCtx)
	go d.serveInLegion(systemCtx)

	// Start Services
	d.startRequiredServices()
	go d.watchServicesDepnedencies(systemCtx)

	// Dump Stats
	startMsg := "I seek to join the Dominion\n" +
		d.config.Dump() +
		"The Dominion ever expands!\n"

	d.Logf(startMsg)

	return d, nil
}

const (
	debugRoutines = "Routine"
	debugRPCs     = "RPC"
	debugLegion   = "Legion"
	debugLocks    = "Locks"
	debugFatal    = "Fatal"
	debugDefault  = "Default"
)

func (d Domain) debugf(class, fmt string, args ...interface{}) {
	const (
		debug       = false
		logRoutines = false
		logRPCs     = false
		logLegion   = true
		logLocks    = false
		logFatal    = true
	)

	if !debug {
		return
	}

	switch class {
	case debugRoutines:
		if logRoutines {
			d.Logf(fmt, args...)
		}
	case debugRPCs:
		if logRPCs {
			d.Logf(fmt, args...)
		}
	case debugLegion:
		if logLegion {
			d.Logf(fmt, args...)
		}
	case debugLocks:
		if logLocks {
			d.Logf(fmt, args...)
		}
	case debugFatal:
		if logFatal {
			d.Logf(fmt, args...)
		}
	case debugDefault:
		d.Logf(fmt, args...)
	default:
		d.Logf("debugf, unknown class:"+fmt, args...)
	}
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
