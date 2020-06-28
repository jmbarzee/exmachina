package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/jmbarzee/domain/server/grpc"
	"github.com/jmbarzee/domain/server/identity"
)

// grabPBIMultiple is used
func (d *Domain) generatePeersPBI() []*pb.Identity {
	pbIdentities := make([]*pb.Identity, 0)

	d.peerMap.Range(func(uuid string, peer *Peer) bool {
		var pbIdent *pb.Identity
		d.debugf(debugLocks, "grabPBIMultiple() pre-lock(%v)\n", peer.UUID)
		peer.RLock()
		{
			d.debugf(debugLocks, "grabPBIMultiple() in-lock(%v)\n", peer.UUID)
			var err error
			pbIdent, err = identity.ConvertItoPBI(peer.Identity)
			if err != nil {
				goto Unlock
			}

			pbIdentities = append(pbIdentities, pbIdent)

		Unlock:
		}
		peer.RUnlock()
		d.debugf(debugLocks, "grabPBIMultiple() post-lock(%v)\n", peer.UUID)

		return true
	})

	return pbIdentities
}

func (d *Domain) generatePBI() *pb.Identity {
	ip, err := d.config.IP.MarshalText()
	if err != nil {
		d.Panic(errors.New("couldn't marshal IP of self"))
	}

	pbIdent := &pb.Identity{
		UUID:        d.config.UUID,
		Version:     d.config.Version.String(),
		Services:    make([]*pb.ServiceIdentity, 0),
		LastContact: time.Now().UnixNano(),
		IP:          ip,
		Port:        int32(d.config.Port),
	}

	d.debugf(debugLocks, "generatePBI() pre-lock(%v)\n", "servicesLock")
	d.servicesLock.Lock()
	{
		d.debugf(debugLocks, "generatePBI() in-lock(%v)\n", "servicesLock")
		for name, service := range d.services {
			service := &pb.ServiceIdentity{
				Name:        name,
				Port:        int32(service.ServiceIdentity.Port),
				LastContact: service.ServiceIdentity.LastContact.UnixNano()}
			pbIdent.Services = append(pbIdent.Services, service)
		}
	}
	d.servicesLock.Unlock()
	d.debugf(debugLocks, "generatePBI() post-lock(%v)\n", "servicesLock")

	return pbIdent
}

func (d *Domain) updateIdentities(identities []identity.Identity) error {
	d.debugf(debugDomain, "updateIdentities()\n")

	for _, identity := range identities {

		if identity.UUID == d.config.UUID {
			// don't add self
			continue
		}

		err := d.updateIdentity(identity)
		if err != nil {
			d.Logf("failed to update identity: %v", err)
		}
	}

	d.debugf(debugDomain, "updateIdentities()\n")
	return nil
}

func (d *Domain) updateIdentity(identity identity.Identity) error {
	d.debugf(debugDomain, "updateIdentity()\n")

	// check if we have peer already
	peer, ok := d.peerMap.Load(identity.UUID)
	if !ok {
		// identity is new
		d.Logf("found new peer: %v\n", identity.UUID)
		newPeer := newPeer(identity)
		d.peerMap.Store(newPeer.UUID, newPeer)
		return nil
	}

	// peer should be garunteed to exist at this point
	peer, ok = d.peerMap.Load(identity.UUID)
	if !ok {
		err := fmt.Errorf("updateIdentity() - failed to find peer, should have existed. %v\n", identity.UUID)
		d.debugf(debugDomain, err.Error())
		return err
	}

	d.debugf(debugLocks, "updateIdentity() pre-lock(%v)\n", peer.UUID)
	peer.Lock()
	{
		d.debugf(debugLocks, "updateIdentity() in-lock(%v)\n", peer.UUID)
		if time.Since(peer.LastContact) > time.Since(identity.LastContact) {
			// found more recent contact with peer
			peer.LastContact = identity.LastContact

			if !peer.IP.Equal(identity.IP) && identity.IP != nil {
				// TODO consider just closing connection instead of restarting it
				d.Logf("%v changed IP form %v to %v\n", peer.UUID, peer.IP, identity.IP)
				peer.IP = identity.IP
				ctx, cancel := context.WithTimeout(context.Background(), d.config.ConnectionConfig.DialTimeout)
				defer cancel()
				err := peer.reconnect(ctx)
				if err != nil {
					// don't fail, look at the rest of the identities.
					d.Logf("failed to reconnect with %v at %s\n", peer.UUID, peer.addr())
					goto Unlock
				}
			}

			if peer.Port != identity.Port && identity.Port != 0 {
				d.Logf("%v changed IP form %v to %v\n", peer.UUID, peer.IP, identity.IP)
				peer.IP = identity.IP
				ctx, cancel := context.WithTimeout(context.Background(), d.config.ConnectionConfig.DialTimeout)
				defer cancel()
				err := peer.reconnect(ctx)
				if err != nil {
					// don't fail, look at the rest of the identities.
					d.Logf("failed to reconnect with %v at %s\n", peer.UUID, peer.addr())
					goto Unlock
				}
			}

			if peer.Version.Compare(identity.Version) < 0 {
				d.Logf("%v changed version form %v to %v\n", peer.UUID, peer.Version.String(), identity.Version.String())
				peer.Version = identity.Version
			}

			for serviceName, serviceIdent := range identity.Services {
				if knownServiceIdent, ok := peer.Services[serviceName]; ok {
					if knownServiceIdent.Port != serviceIdent.Port {
						d.Logf("%v moved \"%v\" from %v to %v\n", peer.UUID, serviceName, knownServiceIdent.Port, serviceIdent.Port)
					}
					peer.Services[serviceName] = serviceIdent

				} else {
					d.Logf("%v started a new service \"%v\"\n", peer.UUID, serviceName)
					peer.Services[serviceName] = serviceIdent
				}
			}

		}
	}

Unlock:
	peer.Unlock()
	d.debugf(debugLocks, "updateIdentity() post-lock(%v)\n", peer.UUID)
	return nil
}
