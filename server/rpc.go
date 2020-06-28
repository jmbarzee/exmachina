package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/jmbarzee/domain/server/grpc"
	"github.com/jmbarzee/domain/server/identity"
)

// GetServices implements grpc and allows the domains to use grpc.
// GetServices serves as the directory of services hosted on all domains.
// GetServices is called by services hosted on a single domain to find their dependencies.
func (d *Domain) GetServices(ctx context.Context, request *pb.GetServicesRequest) (*pb.GetServicesReply, error) {
	d.Logf("[RPC] GetServices")
	serviceName := request.Name
	addrs := d.findService(serviceName)
	reply := &pb.GetServicesReply{
		Addresses: addrs,
	}
	return reply, nil
}

// DumpIdentityList implements grpc and allows the domain to use grpc.
// DumpIdentityList serves as the heartbeat between domains.
func (d *Domain) DumpIdentityList(ctx context.Context, request *pb.DumpIdentityListRequest) (*pb.IdentityListReply, error) {
	d.Logf("[RPC] DumpIdentityList")

	// Prepare reply
	reply := &pb.IdentityListReply{
		Identity:     d.generatePBI(),
		IdentityList: d.generatePeersPBI(),
	}

	d.debugf(debugRPCs, "DumpIdentityList(ctx) returning\n")
	return reply, nil
}

// ShareIdentityList implements grpc and allows the domain to use grpc.
// ShareIdentityList serves as the heartbeat between domains.
func (d *Domain) ShareIdentityList(ctx context.Context, request *pb.IdentityListRequest) (*pb.IdentityListReply, error) {
	d.Logf("[RPC] ShareIdentityList from %v", request.GetIdentity().GetUUID())

	// Parse request
	ident, err := identity.ConvertPBItoI(request.GetIdentity())
	if err != nil {
		d.Logf("Failed to parse identity from request: %v", err.Error())
		return nil, err
	}
	ident.LastContact = time.Now()
	err = d.updateIdentity(ident)
	if err != nil {
		d.Logf("rpc failed to update identity of sender: %v", err)
		return nil, err
	}

	identities, err := identity.ConvertPBItoIMultiple(request.GetIdentityList())
	if err != nil {
		d.Logf("rpc failed to convert identities of sender's peers: %v", err)
		return nil, err
	}

	// Handle RPC
	err = d.updateIdentities(identities)
	if err != nil {
		d.Logf("Couldn't update Identities: %v", err.Error())
	}

	// Prepare reply
	reply := &pb.IdentityListReply{
		Identity:     d.generatePBI(),
		IdentityList: d.generatePeersPBI(),
	}

	d.debugf(debugRPCs, "ShareIdentityList(ctx, %v) returning\n", request.GetIdentity().GetUUID())
	return reply, nil
}

// rpcShareIdentityList calls the grpc ShareIdentityList on the provided peer.
func (d *Domain) rpcShareIdentityList(ctx context.Context, peer *Peer) error {
	// d.debugf(debugRPCs, "rpcShareIdentityList(%v)\n", peer.UUID)
	err := d.checkConnection(peer)
	if err != nil {
		d.Logf("failed to checkConnection(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	var reply *pb.IdentityListReply

	d.debugf(debugLocks, "rpcShareIdentityList() pre-lock(%v)\n", peer.UUID)
	peer.RLock()
	{
		d.debugf(debugLocks, "rpcShareIdentityList() in-lock(%v)\n", peer.UUID)

		// Prepare request
		request := &pb.IdentityListRequest{
			Identity:     d.generatePBI(),
			IdentityList: d.generatePeersPBI(),
		}

		// Send RPC
		// d.Logf("rpcShareIdentityList   -> uuid:%v %v\n", peer.UUID, peer.addr())
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
	identities, err := identity.ConvertPBItoIMultiple(reply.GetIdentityList())
	if err != nil {
		return err
	}

	ident, err := identity.ConvertPBItoI(reply.GetIdentity())
	if err != nil {
		d.Logf(err.Error())
	} else {
		identities = append(identities, ident)
	}
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
	d.Logf("[RPC] OpenPosition from %v", request.GetIdentity().GetUUID())

	reply := &pb.OpenPositionReply{}
	var err error

	d.debugf(debugLocks, "OpenPosition() pre-lock()\n")
	d.electionsLock.Lock()
	{
		d.debugf(debugLocks, "OpenPosition() in-lock()\n")

		if election, ok := d.elections[request.GetName()]; ok {
			d.debugf(debugDefault, "  rejecting, found a pending election: %+v\n", election)
			reply.Accept = false
		} else {
			election := Election{
				Start:   time.Now(),
				SelfRun: false,
			} // Just acts as a tracking device
			d.beginElection(context.Background(), request.GetName(), &election)
			serviceConfig, err := d.serviceConfigFromName(request.GetName())
			if err != nil {
				d.Logf("election proposed for unknown service: %s\n", request.GetName())
			} else {
				reply.Proficiency = d.getProficiencyForService(serviceConfig)
				reply.Accept = true
			}
		}
	}
	d.electionsLock.Unlock()
	d.debugf(debugLocks, "OpenPosition() post-lock()\n")
	return reply, err
}

// rpcOpenPosition calls the grpc OpenPosition on the provided peer.
func (d *Domain) rpcOpenPosition(ctx context.Context, peer *Peer, serviceName string, ballots chan<- Ballot) error {
	d.debugf(debugRPCs, "rpcOpenPosition(%v)\n", peer.UUID)
	err := d.checkConnection(peer)
	if err != nil {
		d.Logf("failed to checkConnection(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	d.debugf(debugLocks, "rpcOpenPosition() pre-lock(%v)\n", peer.UUID)
	peer.RLock()
	{

		request := &pb.OpenPositionRequest{
			Identity: d.generatePBI(),
			Name:     serviceName,
		}

		// Send RPC
		// d.Logf("rpcOpenPosition   -> uuid:%v %v\n", peer.UUID, peer.addr())
		client := pb.NewDomainClient(peer.conn)
		reply, err := client.OpenPosition(ctx, request)
		if err != nil {
			d.debugf(debugDefault, "  ballot had error: %v\n", err)
		} else {
			peer.LastContact = time.Now()
			ballots <- Ballot{
				Accept:      reply.GetAccept(),
				Proficiency: reply.GetProficiency(),
				UUID:        peer.UUID,
			}
		}

	}
	peer.RUnlock()
	d.debugf(debugLocks, "rpcOpenPosition() post-lock(%v)\n", peer.UUID)

	if err != nil {
		d.Logf("failed to rpcOpenPosition(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	d.debugf(debugRPCs, "rpcOpenPosition(%v) returning\n", peer.UUID)
	return nil
}

// ClosePosition implements grpc and allows the domains to use grpc.
// ClosePosition serves as the begining of an election for domains.
func (d *Domain) ClosePosition(ctx context.Context, request *pb.ClosePositionRequest) (*pb.ClosePositionReply, error) {
	d.Logf("[RPC] ClosePosition from %v", request.GetIdentity().GetUUID())

	reply := &pb.ClosePositionReply{}
	var err error

	d.debugf(debugLocks, "ClosePosition() pre-lock()\n")
	d.electionsLock.Lock()
	{
		d.debugf(debugLocks, "ClosePosition() in-lock()\n")

		if _, ok := d.elections[request.GetName()]; ok {
			delete(d.elections, request.GetName())
			if request.GetElected() {
				serviceConfig, err := d.serviceConfigFromName(request.GetName())
				if err == nil {
					err = d.startService(serviceConfig)
					if err != nil {
						d.Logf("Failed to start service after acceptance: %v", err)
					} else {
						reply.Accept = true
					}
				}
			}
		} else if request.GetElected() {
			err = errors.New("Was elected with no knowledge of election process! did the election expire?")
			d.Logf("%v", err)
		}
	}
	d.electionsLock.Unlock()
	d.debugf(debugLocks, "ClosePosition() post-lock()\n")

	return reply, err
}

// rpcClosePosition calls the grpc ClosePosition on the provided peer.
func (d *Domain) rpcClosePosition(ctx context.Context, peer *Peer, serviceName string, elected bool) error {
	d.debugf(debugRPCs, "rpcClosePosition(%v)\n", peer.UUID)
	err := d.checkConnection(peer)
	if err != nil {
		d.Logf("failed to rpcClosePosition(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	d.debugf(debugLocks, "rpcClosePosition() pre-lock(%v)\n", peer.UUID)
	peer.RLock()
	{
		request := &pb.ClosePositionRequest{
			Identity: d.generatePBI(),
			Name:     serviceName,
			Elected:  elected,
		}

		// Send RPC
		// d.Logf("rpcClosePosition   -> uuid:%v %v\n", peer.UUID, peer.addr())
		client := pb.NewDomainClient(peer.conn)
		_, err = client.ClosePosition(ctx, request)
		if err != nil {
			peer.LastContact = time.Now()
			// TODO @jmbarzee check reply from ClosePosition
		}

	}
	peer.RUnlock()
	d.debugf(debugLocks, "rpcClosePosition() post-lock(%v)\n", peer.UUID)

	if err != nil {
		d.Logf("failed to rpcClosePosition(%v) - %v\n", peer.UUID, err.Error())
		return err
	}

	d.debugf(debugRPCs, "rpcClosePosition(%v) returning\n", peer.UUID)
	return nil
}
