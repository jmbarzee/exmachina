package server

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
		if !d.hasTraitsForService(serviceConfig) {
			continue
		}

		// Domain CAN and MUST start the service
		err := d.startService(serviceConfig)
		if err != nil {
			d.Logf("startRequiredServices(): startService() failed: %v", err)
		}
	}
}

func (d *Domain) watchServicesDepnedencies(ctx context.Context) {
	d.debugf(debugRoutines, "watchServicesDepnedencies()\n")

	ticker := time.NewTicker(d.config.ServiceHierarchyConfig.DependencyCheck.Get())

Loop:
	for {
		select {
		case <-ticker.C:
			dependencies := make(map[string]int)

			// Collect all dependencies
			d.debugf(debugLocks, "watchServicesDepnedencies() pre-lock(%v)\n", "servicesLock")
			d.servicesLock.Lock()
			{
				d.debugf(debugLocks, "watchServicesDepnedencies() in-lock(%v)\n", "servicesLock")
				for _, service := range d.services {
					for _, dependency := range service.ServiceConfig.Depends {
						dependencies[dependency]++
					}
				}
			}
			d.servicesLock.Unlock()
			d.debugf(debugLocks, "watchServicesDepnedencies() post-lock(%v)\n", "servicesLock")

			// Check that dependencies exist
			for dependency := range dependencies {
				services := d.findService(dependency)
				if len(services) == 0 {
					if _, ok := d.elections[dependency]; !ok {
						go d.hostElection(ctx, dependency)
					}
				}
			}

		case <-ctx.Done():
			break Loop
		}
	}
	d.debugf(debugRoutines, "watchServicesDepnedencies() stopping\n")
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

	d.debugf(debugLocks, "watchServicesDepnedencies() pre-lock(%v)\n", "servicesLock")
	d.servicesLock.Lock()
	{
		d.debugf(debugLocks, "watchServicesDepnedencies() in-lock(%v)\n", "servicesLock")
		for ownedServiceName := range d.services {
			if serviceName == ownedServiceName {
				addr := fmt.Sprintf("%s:%v", d.config.IP.String(), d.config.Port)
				serviceAddrs = append(serviceAddrs, addr)
			}
		}
	}
	d.servicesLock.Unlock()
	d.debugf(debugLocks, "watchServicesDepnedencies() post-lock(%v)\n", "servicesLock")

	d.peerMap.Range(func(uuid string, peer *Peer) bool {
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
