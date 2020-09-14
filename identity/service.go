package identity

import pb "github.com/jmbarzee/dominion/grpc"

// ServiceIdentity represents a service running under a domain
type ServiceIdentity struct {
	// Type is the type of the service
	Type string
	// Address is the network address of the service
	Address Address
}

func (i ServiceIdentity) String() string {
	return `{ Address: ` + i.Address.String() + `, Type: ` + i.Type + `}`
}

// NewServiceIdentity creates a ServiceIdentity from a pb.ServiceIdentity
func NewServiceIdentity(pbsIdent *pb.ServiceIdentity) ServiceIdentity {
	return ServiceIdentity{
		Type:    pbsIdent.GetType(),
		Address: NewAddress(pbsIdent.GetAddress()),
	}
}

// NewPBServiceIdentity creates a pb.ServiceIdentity from a ServiceIdentity
func NewPBServiceIdentity(sIdent ServiceIdentity) *pb.ServiceIdentity {
	return &pb.ServiceIdentity{
		Type:    sIdent.Type,
		Address: NewPBAddress(sIdent.Address),
	}
}

// NewServiceIdentityMap creates a map of new ServiceIdentitys from a map of pb.ServiceIdentity
func NewServiceIdentityMap(pbsIdents map[string]*pb.ServiceIdentity) map[string]ServiceIdentity {
	sIdents := make(map[string]ServiceIdentity, len(pbsIdents))
	for _, pbsIdent := range pbsIdents {
		sIdents[pbsIdent.GetType()] = NewServiceIdentity(pbsIdent)
	}
	return sIdents
}

// NewPBServiceIdentityMap creates a map of new ServiceIdentitys from a map of pb.ServiceIdentity
func NewPBServiceIdentityMap(sIdents map[string]ServiceIdentity) map[string]*pb.ServiceIdentity {
	pbsIdents := make(map[string]*pb.ServiceIdentity, len(sIdents))
	for _, sIdent := range sIdents {
		pbsIdents[sIdent.Type] = NewPBServiceIdentity(sIdent)
	}
	return pbsIdents
}

// NewServiceIdentityList creates a list of new ServiceIdentitys from a list of pb.ServiceIdentity
func NewServiceIdentityList(pbsIdents []*pb.ServiceIdentity) []ServiceIdentity {
	sIdents := make([]ServiceIdentity, len(pbsIdents))
	for _, pbsIdent := range pbsIdents {
		sIdents = append(sIdents, NewServiceIdentity(pbsIdent))
	}
	return sIdents
}

// NewPBServiceIdentityList creates a list of new ServiceIdentitys from a list of pb.ServiceIdentity
func NewPBServiceIdentityList(sIdents []ServiceIdentity) []*pb.ServiceIdentity {
	pbsIdents := make([]*pb.ServiceIdentity, len(sIdents))
	for _, sIdent := range sIdents {
		pbsIdents = append(pbsIdents, NewPBServiceIdentity(sIdent))
	}
	return pbsIdents
}
