package identity

import (
	"fmt"

	"github.com/blang/semver"
	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/system"
)

type (
	// DominionIdentity represents a dominion
	DominionIdentity struct {
		// Address is the network address of the dominion
		Address Address
		// Version is the version of code being run by the dominion
		Version semver.Version
	}
)

func (i DominionIdentity) String() string {
	return `{
	Address: ` + i.Address.String() + `
	Version: ` + i.Version.String() + `
}`
}

// NewDominionIdentity creates a DominionIdentity from a pb.DominionIdentity
func NewDominionIdentity(pbDsIdent *pb.DominionIdentity) DominionIdentity {

	version, err := semver.Parse(pbDsIdent.GetVersion())
	if err != nil {
		system.Panic(fmt.Errorf("Error parseing version from \"%v\" - %v", pbDsIdent.GetVersion(), err.Error()))
	}

	return DominionIdentity{
		Address: NewAddress(pbDsIdent.GetAddress()),
		Version: version,
	}
}

// NewPBDominionIdentity creates a pb.DominionIdentity from a DominionIdentity
func NewPBDominionIdentity(DsIdent DominionIdentity) *pb.DominionIdentity {
	return &pb.DominionIdentity{
		Address: NewPBAddress(DsIdent.Address),
		Version: DsIdent.Version.String(),
	}
}
