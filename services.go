package domain

import (
	"context"
	"fmt"
	"time"
)

func (d *Domain) startRequiredServices() {
	d.Logf("Starting required services")
	for _, serviceConfig := range d.config.Services {
		if serviceConfig.Priority != Required {
			continue
		}

		// Check to see if domain posses the nessecarry traits
		hasAllTraits := true
		for _, trait := range serviceConfig.Traits {
			hasAllTraits = hasAllTraits && d.hasTrait(trait)
		}
		if !hasAllTraits {
			continue
		}

		// Domain CAN and MUST start the service
		d.startService(serviceConfig)
	}
}

func (d *Domain) watchServicesDepnedencies(ctx context.Context) {
	d.debugf(debugRoutines, "watchServicesDepnedencies()\n")
	d.Logf("Starting required services")

	ticker := time.NewTicker(d.config.TimingConfig.ServiceDependsCheck.Get())

Loop:
	for {
		select {
		case <-ticker.C:
			dependencies := make(map[string]int)

			// Collect all dependencies
			d.debugf(debugLocks, "watchServicesDepnedencies() pre-lock(%v)\n", "ServicesLock")
			d.ServicesLock.Lock()
			{
				d.debugf(debugLocks, "watchServicesDepnedencies() in-lock(%v)\n", "ServicesLock")
				for _, service := range d.services {
					for _, dependency := range service.ServiceConfig.Depends {
						dependencies[dependency]++
					}
				}
			}
			d.ServicesLock.Unlock()
			d.debugf(debugLocks, "watchServicesDepnedencies() post-lock(%v)\n", "ServicesLock")

			// Check that dependencies exist
			for dependency := range dependencies {
				services := d.findService(dependency)
				if len(services) == 0 {
					go d.holdElection(ctx, dependency)
				}
			}

		case <-ctx.Done():
			break Loop
		}
	}
	d.debugf(debugRoutines, "watchServicesDepnedencies() stopping\n")
}

func (d *Domain) holdElection(ctx context.Context, serviceName string) {
	d.Logf("Holding Election for: %s", serviceName)
}

func (d *Domain) hasTrait(trait string) bool {
	for _, ownTrait := range d.config.Traits {
		if ownTrait == trait {
			return true
		}
	}
	return false
}

func (d *Domain) findService(serviceName string) []string {
	serviceAddrs := make([]string, 0)

	d.debugf(debugLocks, "watchServicesDepnedencies() pre-lock(%v)\n", "ServicesLock")
	d.ServicesLock.Lock()
	{
		d.debugf(debugLocks, "watchServicesDepnedencies() in-lock(%v)\n", "ServicesLock")
		for ownedServiceName := range d.services {
			if serviceName == ownedServiceName {
				addr := fmt.Sprintf("%s:%v", d.config.IP.String(), d.config.Port)
				serviceAddrs = append(serviceAddrs, addr)
			}
		}
	}
	d.ServicesLock.Unlock()
	d.debugf(debugLocks, "watchServicesDepnedencies() post-lock(%v)\n", "ServicesLock")

	d.peerMap.Range(func(uuid string, peer *peer) bool {
		d.debugf(debugLocks, "ShareIdentityList() pre-lock(%v)\n", peer.UUID)
		peer.RLock()
		{
			d.debugf(debugLocks, "ShareIdentityList() in-lock(%v)\n", peer.UUID)
			for peerServiceName, port := range peer.Services {
				if serviceName == peerServiceName {
					addr := fmt.Sprintf("%s:%v", peer.IP.String(), port)
					serviceAddrs = append(serviceAddrs, addr)
				}
			}
		}
		peer.RUnlock()
		d.debugf(debugLocks, "updateLegion() post-lock(%v)\n", peer.UUID)

		return true
	})

	return serviceAddrs
}
