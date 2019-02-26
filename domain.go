package domain

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
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

		// services is the list of services the domain currently offers
		services     map[string]Service
		ServicesLock sync.Mutex

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

		// TimingConfig specifies all the timings needed for autoconnecting and heartbeats
		TimingConfig TimingConfig

		// Log is where the logs from the domain are left.
		Log *log.Logger
	}

	TimingConfig struct {
		// DialTimeout is how long a domain will wait for a grpc.ClientConn to establish
		DialTimeout time.Duration

		// IsolationCheck is the range of possible durations between isolation checks
		IsolationCheck RangeTiming
		// IsolationTimeout is the range of possible durations after which a domain will determine it is isolated
		IsolationTimeout RangeTiming

		// HeartbeatCheck is the range of possible durations after which a domain will send a heartbeat
		HeartbeatCheck RangeTiming
	}

	RangeTiming struct {
		// DialTimeout is the top of the possible durations
		Upper time.Duration
		// DialTimeout is the bottom of the possible durations
		Lower time.Duration
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

	// Dump Stats
	startMsg := "I seek to join the Dominion\n" +
		d.config.Dump() +
		"\n" +
		"The Dominion ever expands!\n" +
		"Long grow the dominion!\n"

	d.Logf(startMsg)

	return d, nil
}

func (c DomainConfig) Check() error {
	if c.UUID == "" {
		return fmt.Errorf("UUID was not set")
	}
	if c.Title == "" {
		return fmt.Errorf("Title was not set")
	}

	if len(c.Traits) == 0 {
		return fmt.Errorf("Traits were not set")
	}
	if len(c.Services) == 0 {
		return fmt.Errorf("Services were not set")
	}

	if c.Port == 0 {
		return fmt.Errorf("Port was not set")
	}

	err := c.TimingConfig.Check()
	if err != nil {
		return err
	}

	if c.Log == nil {
		return fmt.Errorf("LogFileName was not set")
	}
	return nil
}

func (c DomainConfig) Dump() string {

	dumpMsg := "\tUUID: " + c.UUID + "\n" +
		"\tTitle: " + c.Title + "\n" +
		"\tVersion: " + c.Version.String() + "\n" +
		"\tTraits: \n"

	dumpMsg += "\tTraits: ["
	for _, trait := range c.Traits {
		dumpMsg += trait + ", "
	}
	dumpMsg += "]\n"

	dumpMsg += "\tServices: [\n"
	for _, serviceConfig := range c.Services {
		dumpMsg += "\t\t" + serviceConfig.String() + ", \n"
	}
	dumpMsg += "\t]\n"

	dumpMsg += "\tAddress: " + c.IP.String() + ":" + strconv.Itoa(c.Port) + "\n" +
		"\tTimingConfig: \n" +
		c.TimingConfig.String()
	return dumpMsg
}

func (c TimingConfig) Check() error {
	if c.DialTimeout == 0 {
		return fmt.Errorf("TimingConfig.DialTimeout was not set")
	}
	if c.IsolationCheck.Upper == 0 {
		return fmt.Errorf("TimingConfig.IsolationCheck.Upper was not set")
	}
	if c.IsolationCheck.Lower == 0 {
		return fmt.Errorf("TimingConfig.IsolationCheck.Lower was not set")
	}
	if c.IsolationTimeout.Upper == 0 {
		return fmt.Errorf("TimingConfig.IsolationTimeout.Upper was not set")
	}
	if c.IsolationTimeout.Lower == 0 {
		return fmt.Errorf("TimingConfig.IsolationTimeout.Lower was not set")
	}
	if c.HeartbeatCheck.Upper == 0 {
		return fmt.Errorf("TimingConfig.HeartbeatCheck.Upper was not set")
	}
	if c.HeartbeatCheck.Lower == 0 {
		return fmt.Errorf("TimingConfig.HeartbeatCheck.Lower was not set")
	}
	return nil
}

func (c TimingConfig) String() string {
	dumpMsg := "\t\tDialTimeout: " + c.DialTimeout.String() + "\n" +
		"\t\tIsolationCheck: " + c.IsolationCheck.String() + "\n" +
		"\t\tIsolationTimeout: " + c.IsolationTimeout.String() + "\n" +
		"\t\tHeartbeatCheck: " + c.HeartbeatCheck.String() + "\n"
	return dumpMsg
}

func (r RangeTiming) String() string {
	return "(" + r.Lower.String() + "," + r.Upper.String() + ")"
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
