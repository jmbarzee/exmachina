package identity

import (
	"fmt"
	"net"
	"time"

	"github.com/blang/semver"
	pb "github.com/jmbarzee/domain/server/grpc"
)

type (
	// Identity contains the all shareable information about a domain
	Identity struct {
		// UUID is a unique identifier for a Domain
		UUID string
		// Version is the version of Code which the Domain is running
		Version semver.Version
		// Services is the list of services the Domain currently offers
		Services map[string]ServiceIdentity

		// LastContact is when the legion last heard from this Identity
		LastContact time.Time

		// IP is the port which the Domain will be responding on
		IP net.IP
		// Port is the port which the Domain will be responding on
		Port int
	}

	ServiceIdentity struct {
		Port        int
		LastContact time.Time
	}
)

func ConvertPBItoIMultiple(pbIdents []*pb.Identity) ([]Identity, error) {
	identities := make([]Identity, 0)

	for _, pbIdent := range pbIdents {
		ident, err := ConvertPBItoI(pbIdent)
		if err != nil {
			return nil, err
			continue
		}
		identities = append(identities, ident)
	}

	return identities, nil
}
func ConvertPBItoI(pbIdent *pb.Identity) (Identity, error) {

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

func ConvertItoPBI(ident Identity) (*pb.Identity, error) {
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
