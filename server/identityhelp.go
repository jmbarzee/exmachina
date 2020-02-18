package server

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/blang/semver"
	pb "github.com/jmbarzee/domain/server/grpc"
)

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

func (d *Domain) convertPBItoIMultiple(pbIdents []*pb.Identity) []Identity {
	identities := make([]Identity, 0)

	for _, pbIdent := range pbIdents {
		ident, err := convertPBItoI(pbIdent)
		if err != nil {
			d.Logf(err.Error())
			continue
		}
		identities = append(identities, ident)
	}

	return identities
}
func convertPBItoI(pbIdent *pb.Identity) (Identity, error) {

	version, err := semver.Parse(pbIdent.GetVersion())
	if err != nil {
		return Identity{}, fmt.Errorf("Error parseing version from \"%v\" - %v", pbIdent.GetVersion(), err.Error())
	}

	ident := Identity{
		UUID:        pbIdent.GetUUID(),
		Version:     version,
		Services:    make(map[string]ServiceIdentity),
		LastContact: time.Unix(0, pbIdent.GetLastContact()),
		IP:          net.ParseIP(string(pbIdent.GetIP())),
		Port:        int(pbIdent.GetPort()),
	}

	for _, service := range pbIdent.GetServices() {
		ident.Services[service.GetName()] = ServiceIdentity{
			Port:        int(pbIdent.GetPort()),
			LastContact: time.Unix(0, service.GetLastContact()),
		}
	}
	return ident, nil
}

// grabPBIMultiple is used
func (d *Domain) grabPBIMultiple() []*pb.Identity {
	pbIdentities := make([]*pb.Identity, 0)

	d.peerMap.Range(func(uuid string, peer *Peer) bool {
		var pbIdent *pb.Identity
		d.debugf(debugLocks, "ShareIdentityList() pre-lock(%v)\n", peer.UUID)
		peer.RLock()
		{
			d.debugf(debugLocks, "ShareIdentityList() in-lock(%v)\n", peer.UUID)
			var err error
			pbIdent, err = convertItoPBI(peer.Identity)
			if err != nil {
				goto Unlock
			}

			pbIdentities = append(pbIdentities, pbIdent)

		Unlock:
		}
		peer.RUnlock()
		d.debugf(debugLocks, "updateLegion() post-lock(%v)\n", peer.UUID)

		return true
	})

	return pbIdentities
}
func convertItoPBI(ident Identity) (*pb.Identity, error) {
	ip, err := ident.IP.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal IP of %v - %v\n", ident.UUID, err)
	}

	pbIdent := &pb.Identity{
		UUID:        ident.UUID,
		Version:     ident.Version.String(),
		Services:    make([]*pb.ServiceIdentity, 0),
		LastContact: ident.LastContact.UnixNano(),
		IP:          ip,
		Port:        int32(ident.Port),
	}

	for name, serviceIdentity := range ident.Services {
		service := &pb.ServiceIdentity{
			Name:        name,
			Port:        int32(serviceIdentity.Port),
			LastContact: serviceIdentity.LastContact.UnixNano()}
		pbIdent.Services = append(pbIdent.Services, service)
	}

	return pbIdent, nil
}
