package domain

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

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

		// Services is the list of services the domain currently offers
		Services map[string]Service

		// peerMap stores the members of a Dominion in a wrapped sync.map
		peerMap *peerMap

		// config does what it is. See DomainConfig for a clear understanding.
		// It is unchanged after NewDomain() completes
		config DomainConfig
	}

	DomainConfig struct {
		// UUID is a unique identifier for a domain
		UUID string
		// Title is the name of the Dominion which the domain belongs to
		Title string

		// Version is the version of Code which the domain is running
		Version semver.Version

		// Traits is the traits possesed by the domain.
		Traits []string
		// Services is the list of possible services.
		Services map[string]ServiceConfig

		// Port is the port which the domain will be responding on
		Port int
		// IP is the port which the domain will be responding on
		IP net.IP
		// DialTimeout is how long a domain will wait for a grpc.ClientConn to establish
		DialTimeout time.Duration

		// Log is where the logs from the domain are left.
		Log *log.Logger
	}
)

// NewDomain creates a new Domain, to correctly build the Domain, just initilize
func NewDomain(ctx context.Context, config DomainConfig) (*Domain, error) {

	// Check config
	if config.UUID == "" {
		return nil, fmt.Errorf("UUID was not set by ConfigureFunc")
	}
	if config.Title == "" {
		return nil, fmt.Errorf("Title was not set by ConfigureFunc")
	}

	if len(config.Traits) == 0 {
		return nil, fmt.Errorf("Traits were not set by ConfigureFunc")
	}
	if len(config.Services) == 0 {
		return nil, fmt.Errorf("Services were not set by ConfigureFunc")
	}

	if config.Port == 0 {
		return nil, fmt.Errorf("Port was not set by ConfigureFunc")
	}
	if config.DialTimeout == 0 {
		return nil, fmt.Errorf("DialTimeout was not set by ConfigureFunc")
	}

	if config.Log == nil {
		return nil, fmt.Errorf("LogFileName was not set by ConfigureFunc")
	}

	// Initilize Remaining Config vvvv

	var err error
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
	}

	// Initilize System (logging, context, and signals)
	var systemCtx context.Context
	systemCtx, d.System = system.NewSystem(ctx, config.Log)

	// Start Routines
	go d.watchIsolation(systemCtx)
	go d.listenForBroadcasts(systemCtx)
	go d.serveInLegion(systemCtx)

	// Start Services
	// go d.hostServices()

	// Dump Stats
	d.Logf("I seek to join the Dominion\n")
	d.Logf("	UUID:%v\n", d.config.UUID)
	d.Logf("	Title:%v\n", d.config.Title)
	d.Logf("\n")
	d.Logf("	Version:%v\n", d.config.Version)
	d.Logf("\n")
	// TODO @jmbarzee print traints and services
	d.Logf("\n")
	d.Logf("	Address:%v:%v\n", d.config.IP, d.config.Port)
	d.Logf("	DialTimeout:%v\n", d.config.DialTimeout)
	d.Logf("\n")
	d.Logf("The Dominion ever expands!\n")
	d.Logf("Long grow the dominion!\n")

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
