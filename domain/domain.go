package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/jmbarzee/dominion/domain/config"
	"github.com/jmbarzee/dominion/domain/dominion"
	"github.com/jmbarzee/dominion/domain/service"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
)

type (
	Domain struct {
		identity.DomainIdentity

		Dominion *dominion.DominionGuard

		// services stores the members of a Dominion in a wrapped sync.map as
		//     ServiceType -> Service
		services service.ServiceMap

		stopBroadcastSelf context.CancelFunc
	}
)

// NewDomain creates a new Domain, to correctly build the Domain, just initilize
func NewDomain(configFilePath string) (*Domain, error) {

	// Check config
	if err := config.SetupFromTOML(configFilePath); err != nil {
		return nil, err
	}

	if err := system.Setup(config.GetDomainConfig().LogFilePath); err != nil {
		return nil, err
	}

	ident, err := NewDomainIdentity(config.GetDomainConfig())
	if err != nil {
		return nil, err
	}

	return &Domain{
		services:       service.NewServiceMap(),
		DomainIdentity: ident,
	}, nil
}

// NewDomainIdentity creates a new DomainIdentity
func NewDomainIdentity(domainConfig config.DomainConfig) (identity.DomainIdentity, error) {
	// Initialize Version
	version, err := semver.Parse(system.Version)
	if err != nil {
		return identity.DomainIdentity{}, fmt.Errorf("failed to semver.Parse(%v): %v\n", version, err.Error())
	}

	// Initialize IP
	ip, err := system.GetOutboundIP()
	if err != nil {
		return identity.DomainIdentity{}, fmt.Errorf("failed to find Local IP: %v\n", err.Error())
	}

	return identity.DomainIdentity{
		UUID:    domainConfig.UUID,
		Version: version,
		Traits:  domainConfig.Traits,
		Address: identity.Address{
			IP:   ip,
			Port: domainConfig.Port,
		},
	}, nil
}

func (d Domain) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(d.DomainIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	// Start Auto Connecting Routines
	go d.routineCheck(ctx, "checkIsolation", config.GetDomainConfig().IsolationCheck, d.checkIsolation)
	go d.routineCheck(ctx, "checkServices", config.GetDomainConfig().ServiceCheck, d.checkServices)

	return d.hostDomain(ctx)
}

func (d Domain) packageDomainIdentity() identity.DomainIdentity {
	ident := d.DomainIdentity
	ident.Services = make(map[string]identity.ServiceIdentity)
	d.services.Range(func(serviceType string, serviceGuard *service.ServiceGuard) bool {
		serviceGuard.LatchRead(func(service *service.Service) error {
			ident.Services[serviceType] = service.ServiceIdentity
			return nil
		})
		return true
	})
	return ident
}

func (d *Domain) updateDominion(ident identity.DominionIdentity) error {
	if d.Dominion == nil {
		system.Logf("Joining Dominion %v:", ident.Address.String())
		d.Dominion = dominion.NewDominionGuard(ident)
		return nil
	} else {
		return d.Dominion.LatchWrite(func(dominion *dominion.Dominion) error {
			if dominion.Address.String() != ident.Address.String() {
				return fmt.Errorf("Dominion Address doesn't known dominion")
			} else {
				dominion.LastContact = time.Now()
				return nil
			}
		})
	}
}
