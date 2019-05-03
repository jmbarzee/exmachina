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

		// ConnectionConfig specifies all the values needed for autoconnecting and heartbeats
		ConnectionConfig ConnectionConfig
		// ServiceHierarchyConfig specifies all the values needed for service elections
		ServiceHierarchyConfig ServiceHierarchyConfig

		// Log is where the logs from the domain are left.
		Log *log.Logger
	}

	ConnectionConfig struct {
		// DialTimeout is how long a domain will wait for a grpc.ClientConn to establish
		DialTimeout time.Duration

		// IsolationCheck is the range of possible durations between isolation checks
		IsolationCheck RangeTiming
		// IsolationTimeout is the range of possible durations after which a domain will determine it is isolated
		IsolationTimeout RangeTiming

		// HeartbeatCheck is the range of possible durations after which a domain will send a heartbeat
		HeartbeatCheck RangeTiming
	}

	ServiceHierarchyConfig struct {
		// ElectionTimeout is how long a domain will wait for an election to conclude
		ElectionTimeout time.Duration
		// RequiredPercentage is the number of nodes required for an election to move forward and select a nominee
		RequiredVotePercentage float64
		// ElectionTimeout is the range of possible durations after which a domain will cancel a pending election for a service
		ElectionBackoff RangeTiming

		// DependencyCheck is the range of possible durations between service dependency checks
		DependencyCheck RangeTiming
	}

	RangeTiming struct {
		// Max is the top of the possible durations
		Max time.Duration
		// Min is the bottom of the possible durations
		Min time.Duration
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

	err := c.ConnectionConfig.Check()
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
		"\tConnectionConfig: \n" +
		c.ConnectionConfig.String() +
		"\tServiceHierarchyConfig: \n" +
		c.ServiceHierarchyConfig.String()
	return dumpMsg
}

func (c ConnectionConfig) Check() error {
	if c.DialTimeout == 0 {
		return fmt.Errorf("ConnectionConfig.DialTimeout was not set")
	}

	err := c.IsolationCheck.Check()
	if err != nil {
		return fmt.Errorf("ConnectionConfig.IsolationCheck.%v", err)
	}
	err = c.IsolationTimeout.Check()
	if err != nil {
		return fmt.Errorf("ConnectionConfig.IsolationTimeout.%v", err)
	}
	err = c.HeartbeatCheck.Check()
	if err != nil {
		return fmt.Errorf("ConnectionConfig.HeartbeatCheck.%v", err)
	}
	return nil
}

func (c ConnectionConfig) String() string {
	dumpMsg := "\t\tDialTimeout: " + c.DialTimeout.String() + "\n" +
		"\t\tIsolationCheck: " + c.IsolationCheck.String() + "\n" +
		"\t\tIsolationTimeout: " + c.IsolationTimeout.String() + "\n" +
		"\t\tHeartbeatCheck: " + c.HeartbeatCheck.String() + "\n"
	return dumpMsg
}

func (c ServiceHierarchyConfig) Check() error {
	if c.ElectionTimeout == 0 {
		return fmt.Errorf("ServiceHierarchyConfig.ElectionTimeout was not set")
	}
	if c.RequiredVotePercentage == 0 {
		return fmt.Errorf("ServiceHierarchyConfig.RequiredVotePercentage was not set")
	}
	err := c.ElectionBackoff.Check()
	if err != nil {
		return fmt.Errorf("ServiceHierarchyConfig.ElectionBackoff.%v", err)
	}

	err = c.DependencyCheck.Check()
	if err != nil {
		return fmt.Errorf("ServiceHierarchyConfig.DependencyCheck.%v", err)
	}
	return nil
}

func (c ServiceHierarchyConfig) String() string {
	dumpMsg := "\t\tRequiredVotePercentage: " + strconv.FormatFloat(c.RequiredVotePercentage, 'f', 6, 64) + "\n" +
		"\t\tElectionBackoff: " + c.ElectionBackoff.String() + "\n" +
		"\t\tDependencyCheck: " + c.DependencyCheck.String() + "\n"
	return dumpMsg
}

func (r RangeTiming) Check() error {
	if r.Max == 0 {
		return fmt.Errorf("Upper was not set")
	}
	if r.Min == 0 {
		return fmt.Errorf("Lower was not set")
	}
	return nil
}

func (r RangeTiming) String() string {
	return "(" + r.Min.String() + "," + r.Max.String() + ")"
}

func (r RangeTiming) Get() time.Duration {
	intMax := int64(r.Max)
	intMin := int64(r.Min)
	return time.Duration(rand.Int63n(intMax-intMin) + intMin)
}
