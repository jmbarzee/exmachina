package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/jmbarzee/domain/grpc"
)

// GetServices implements grpc and allows the domains to use grpc.
// GetServices serves as the directory of services hosted on all domains.
// GetServices is called by services hosted on a single domain to find their dependencies.
func (d *Domain) GetServices(ctx context.Context, request *pb.GetServicesRequest) (*pb.GetServicesReply, error) {
	return nil, errors.New("UnImplemented!")
}

// rpcGetServices calls the grpc GetServices on the provided peer.
func (d *Domain) rpcGetServices(ctx context.Context, peer *peer) error {
	return errors.New("UnImplemented!")
}

// ShareIdentityList implements grpc and allows the domain to use grpc.
// ShareIdentityList serves as the heartbeat between domains.
func (d *Domain) ShareIdentityList(ctx context.Context, request *pb.IdentityListRequest) (*pb.IdentityListReply, error) {
	d.debugf(debugRPCs, "ShareIdentityList(ctx, %v)\n", request.GetIdentity().GetUUID())

	d.Logf("rpcShareIdentityList <-   uuid:%v\n", request.GetIdentity().GetUUID())

	// Parse request
	identity, err := convertPBItoI(request.GetIdentity())
	if err != nil {
		d.Logf("Failed to parse identity from request: %v", err.Error())
		return nil, err
	}
	identity.LastContact = time.Now()
	err = d.updateIdentity(identity)
	if err != nil {
		d.Logf("rpc failed to update identity of sender: %v", err)
		return nil, err
	}

	identities := d.convertPBItoIMultiple(request.GetIdentityList())

	// Handle RPC
	err = d.updateIdentities(identities)
	if err != nil {
		d.Logf("Couldn't update Identities: %v", err.Error())
	}

	// Prepare reply
	pbIdent, err := d.generatePBI()
	if err != nil {
		d.Panic(fmt.Errorf("Couldn't convert own Identity to pb.Identity: %v", err.Error()))
	}

	reply := &pb.IdentityListReply{
		Identity:     pbIdent,
		IdentityList: d.grabPBIMultiple(),
	}

	d.debugf(debugRPCs, "ShareIdentityList(ctx, %v) returning\n", request.GetIdentity().GetUUID())
	return reply, nil
}

// rpcShareIdentityList calls the grpc ShareIdentityList on the provided peer.
func (d *Domain) rpcShareIdentityList(ctx context.Context, peer *peer) error {
	d.debugf(debugRPCs, "rpcShareIdentityList(%v)\n", peer.UUID)
	err := d.checkConnection(peer)
	if err != nil {
		d.Logf("failed to checkConnection(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	var reply *pb.IdentityListReply

	d.debugf(debugLocks, "rpcShareIdentityList() pre-lock(%v)\n", peer.UUID)
	peer.RLock() // Dirty Lock
	{
		d.debugf(debugLocks, "rpcShareIdentityList() in-lock(%v)\n", peer.UUID)

		// Prepare request
		pbIdent, err := d.generatePBI()
		if err != nil {
			d.Panic(fmt.Errorf("Couldn't convert own Identity to pb.Identity: %v", err.Error()))
		}

		request := &pb.IdentityListRequest{
			Identity:     pbIdent,
			IdentityList: d.grabPBIMultiple(),
		}

		// Send RPC
		d.Logf("rpcShareIdentityList   -> uuid:%v %v\n", peer.UUID, peer.addr())
		client := pb.NewDomainClient(peer.conn)
		reply, err = client.ShareIdentityList(ctx, request)
		if err != nil {
			peer.LastContact = time.Now()
		}
		// err is checked again after lock

	}
	peer.RUnlock()
	d.debugf(debugLocks, "rpcShareIdentityList() post-lock(%v)\n", peer.UUID)

	if err != nil {
		d.Logf("failed to ShareIdentityList(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	// Parse reply
	// TODO handle reply Identity
	identities := d.convertPBItoIMultiple(reply.GetIdentityList())
	err = d.updateIdentities(identities)
	if err != nil {
		d.Logf("rpcShareIdentityList(%v) updateIdentities failed: %v\n", peer.UUID, err)
		return err
	}

	d.debugf(debugRPCs, "rpcShareIdentityList(%v) returning\n", peer.UUID)
	return nil
}

// OpenPosition implements grpc and allows the domains to use grpc.
// OpenPosition serves as the begining of an election for domains.
func (d *Domain) OpenPosition(ctx context.Context, request *pb.OpenPositionRequest) (*pb.OpenPositionReply, error) {
	return nil, errors.New("UnImplemented!")
}

// rpcOpenPosition calls the grpc OpenPosition on the provided peer.
func (d *Domain) rpcOpenPosition(ctx context.Context, peer *peer) error {
	return errors.New("UnImplemented!")
}

// OfferPosition implements grpc and allows the domains to use grpc.
// OfferPosition serves as the begining of an election for domains.
func (d *Domain) OfferPosition(ctx context.Context, request *pb.OfferPositionRequest) (*pb.OfferPositionReply, error) {
	return nil, errors.New("UnImplemented!")
}

// rpcOfferPosition calls the grpc OfferPosition on the provided peer.
func (d *Domain) rpcOfferPosition(ctx context.Context, peer *peer) error {
	return errors.New("UnImplemented!")
}
