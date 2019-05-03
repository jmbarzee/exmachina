package domain

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grandcat/zeroconf"
	pb "github.com/jmbarzee/domain/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (d *Domain) watchIsolation(ctx context.Context) {
	d.debugf(debugRoutines, "watchIsolation()\n")

	ticker := time.NewTicker(d.config.ConnectionConfig.IsolationCheck.Get())

	var stopBroadcastSelf context.CancelFunc
	stopBroadcastSelf = nil

Loop:
	for {
		select {
		case <-ticker.C:
			// review all connections and check that we've heard recently enough
			//d.debugf(debugDomain, "watchIsolation() reviewing connections \n")

			heartbeatTimeout := d.config.ConnectionConfig.HeartbeatCheck.Get()
			isolationTimeout := d.config.ConnectionConfig.IsolationTimeout.Get()

			lonely := true
			d.peerMap.Range(func(uuid string, peer *Peer) bool {
				d.debugf(debugLocks, "watchIsolation() pre-lock(%v)\n", uuid)
				peer.RLock()
				{
					d.debugf(debugLocks, "watchIsolation() in-lock(%v)\n", uuid)
					if time.Since(peer.LastContact) > heartbeatTimeout {
						// its been a while, make sure they are still alive
						go d.rpcShareIdentityList(ctx, peer)
					}

					if time.Since(peer.LastContact) < isolationTimeout {
						// Not lonely yet
						//d.debugf(debugDefault, "watchIsolation() found a recent peer \n")
						lonely = false
					}

				}
				peer.RUnlock()
				d.debugf(debugLocks, "watchIsolation() post-lock(%v)\n", uuid)
				return true
			})

			if lonely && stopBroadcastSelf == nil {
				// Look for a legion
				d.debugf(debugDomain, "watchIsolation() discovered loneliness \n")
				var ctxBroadcast context.Context
				ctxBroadcast, stopBroadcastSelf = context.WithCancel(ctx)

				go d.broadcastSelf(ctxBroadcast)
			} else if !lonely && stopBroadcastSelf != nil {
				// Stop looking for a legion
				d.debugf(debugDomain, "watchIsolation() no longer lonely \n")
				stopBroadcastSelf()
				stopBroadcastSelf = nil

			}
		case <-ctx.Done():
			break Loop
		}
	}
	if stopBroadcastSelf != nil {
		stopBroadcastSelf()
	}

	d.debugf(debugRoutines, "watchIsolation() stopping\n")
}

// broadcastSelf uses zero conf to broadcast to a network.
func (d *Domain) broadcastSelf(ctx context.Context) {
	d.debugf(debugRoutines, "broadcastSelf()\n")

	// setup broadcasting
	server, err := zeroconf.Register(string(d.config.UUID), d.config.Title, "local.", d.config.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		d.debugf(debugFatal, "Failed to broadcast:", err.Error())
		d.Panic(err)
	}
	d.Logf("Started broadcasting .oO \n")

	<-ctx.Done()
	server.Shutdown()
	d.Logf("Stopped broadcasting\n")
}

// listenForBroadcasts listens for other Legionnaires.
func (d *Domain) listenForBroadcasts(ctx context.Context) {
	d.debugf(debugRoutines, "listenForBroadcasts()\n")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		d.debugf(debugFatal, "listenForBroadcasts() Failed to initialize resolver: %v\n", err.Error())
		d.Panic(err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	err = resolver.Browse(ctx, d.config.Title, "local.", entries)
	if err != nil {
		d.debugf(debugFatal, "listenForBroadcasts() Failed to browse: %v\n", err.Error())
		d.Panic(err)
	}

	d.Logf("Listening for broadcasts...\n")

Loop:
	for {
		select {
		case entry, ok := <-entries:
			if !ok {
				// channel closed
				break Loop
			}

			duuid := d.config.UUID
			dinst := entry.Instance
			if dinst == duuid {
				// don't connect to self
				break
			}

			var conn *grpc.ClientConn
			var ip net.IP
			var port int
			// if len(entry.AddrIPv6) > 0 {
			// 	ip = entry.AddrIPv6[0]
			// 	port = entry.Port
			// 	addr := fmt.Sprintf("[%s]:%v", ip.String(), port)
			// 	conn, err = grpc.Dial(addr, grpc.WithInsecure())
			// 	if err != nil {
			// 		d.debugf(debugDefault, "Found ipv6 but connection failed: %v\n", err.Error())
			// 		err = nil // clear err
			// 		// don't return so that we try ipv4 as well
			// 	}
			// }
			if len(entry.AddrIPv4) <= 0 {
				break
			}

			uuid := entry.Instance
			ip = entry.AddrIPv4[0]
			port = entry.Port
			d.Logf("Found broadcast - uuid:%v ip:%v port:%v\n", uuid, ip, port)

			// TODO check if peer is already known and reconnect / heartbeat instead

			addr := fmt.Sprintf("%s:%v", ip.String(), port)
			conn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				d.debugf(debugFatal, "Found ipv4 but connection failed: %v\n", err.Error())
				d.Panic(err)
			}
			d.debugf(debugDefault, "listenForBroadcasts() connected to %v \n", entry.Instance)

			newPeer := &Peer{
				Identity: Identity{
					UUID:        uuid,
					Services:    make(map[string]ServiceIdentity),
					IP:          ip,
					Port:        port,
					LastContact: time.Now(),
				},

				conn: conn,
			}

			// Add the new member
			d.peerMap.Store(uuid, newPeer)

			// Let peer know our state
			go d.rpcShareIdentityList(ctx, newPeer)

		case <-ctx.Done():
			break Loop
		}
	}

	d.debugf(debugRoutines, "listenForOthers() stopping\n")
}

func (d *Domain) serveInLegion(ctx context.Context) {
	d.debugf(debugRoutines, "serveInLegion()\n")
	// TODO use or handle context

	address := fmt.Sprintf("%s:%v", "", d.config.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		d.debugf(debugFatal, "serveInLegion() Failed to listen: %v\n", err)
		d.Panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterDomainServer(server, d)
	// Register reflection service on gRPC server.
	go func() {
		<-ctx.Done()
		server.GracefulStop()
		d.Logf("Stopped grpc server gracefully. ")
	}()

	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		d.debugf(debugFatal, "serveInLegion() Failed to serve: %v\n", err)
		d.Panic(err)
	}

	d.debugf(debugRoutines, "serveInLegion() stopping\n")
}
