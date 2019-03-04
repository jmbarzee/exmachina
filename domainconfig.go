package domain

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/blang/semver"
)

type (
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

		// ServiceDependsCheck is the range of possible durations between service dependency checks
		ServiceDependsCheck RangeTiming
	}

	RangeTiming struct {
		// DialTimeout is the top of the possible durations
		Upper time.Duration
		// DialTimeout is the bottom of the possible durations
		Lower time.Duration
	}
)

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
		"\tVersion: " + c.Version.String() + "\n"

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

	err := c.IsolationCheck.Check()
	if err != nil {
		return fmt.Errorf("TimingConfig.IsolationCheck.%v", err)
	}
	err = c.IsolationTimeout.Check()
	if err != nil {
		return fmt.Errorf("TimingConfig.IsolationTimeout.%v", err)
	}
	err = c.HeartbeatCheck.Check()
	if err != nil {
		return fmt.Errorf("TimingConfig.HeartbeatCheck.%v", err)
	}
	err = c.ServiceDependsCheck.Check()
	if err != nil {
		return fmt.Errorf("TimingConfig.ServiceDependsCheck.%v", err)
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

func (r RangeTiming) Check() error {
	if r.Upper == 0 {
		return fmt.Errorf("Upper was not set")
	}
	if r.Lower == 0 {
		return fmt.Errorf("Lower was not set")
	}
	return nil
}

func (r RangeTiming) String() string {
	return "(" + r.Lower.String() + "," + r.Upper.String() + ")"
}

func (r RangeTiming) Get() time.Duration {
	intMax := int64(r.Upper)
	intMin := int64(r.Lower)
	return time.Duration(rand.Int63n(intMax-intMin) + intMin)
}
