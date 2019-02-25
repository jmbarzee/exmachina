package domain

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/blang/semver"
)

// Identity contains the all shareable information about a legionnaire
type Identity struct {
	// UUID is a unique identifier for a Leggionnair
	UUID string
	// Version is the version of Code which the Legioonnaire is running
	Version semver.Version
	// Services is the list of services the Legionnaire currently offers
	Services map[string]int

	// LastContact is when the legion last heard from this Identity
	LastContact time.Time

	// IP is the port which the Legionnaire will be responding on
	IP net.IP
	// Port is the port which the Legionnaire will be responding on
	Port int
}

func (d *Domain) updateIdentities(identities []Identity) error {
	d.debugf(debugLegion, "updateIdentities()\n")

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

	d.debugf(debugLegion, "updateIdentities()\n")
	return nil
}

func (d *Domain) updateIdentity(identity Identity) error {
	d.debugf(debugLegion, "updateIdentity()\n")

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
		d.debugf(debugLegion, err.Error())
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
				ctx, cancel := context.WithTimeout(context.Background(), d.config.DialTimeout)
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
				ctx, cancel := context.WithTimeout(context.Background(), d.config.DialTimeout)
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

			// TODO handle service changes

		}
	}

Unlock:
	peer.Unlock()
	d.debugf(debugLocks, "updateIdentity() post-lock(%v)\n", peer.UUID)

	d.debugf(debugLegion, "updateIdentity()\n")
	return nil
}
