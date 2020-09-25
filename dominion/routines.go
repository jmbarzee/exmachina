package dominion

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/jmbarzee/dominion/dominion/config"
	"github.com/jmbarzee/dominion/dominion/domain"
	pb "github.com/jmbarzee/dominion/grpc"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// listenForBroadcasts listens for other Legionnaires.
func (d *Dominion) listenForBroadcasts(ctx context.Context) {
	routineName := "listenForBroadcasts"
	system.LogRoutinef(routineName, "Starting routine")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		system.LogRoutinef(routineName, "Failed to initialize resolver: %v", err.Error())
		system.Panic(err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	err = resolver.Browse(ctx, "dominion", "local.", entries)
	if err != nil {
		system.LogRoutinef(routineName, "Failed to browse: %v", err.Error())
		system.Panic(err)
	}

	system.LogRoutinef(routineName, "Listening for broadcasts...")

Loop:
	for {
		select {
		case entry, ok := <-entries:
			if !ok {
				// channel closed
				break Loop
			}
			if len(entry.AddrIPv4) <= 0 {
				break
			}

			uuid := entry.Instance
			ip := entry.AddrIPv4[0]
			port := entry.Port
			system.LogRoutinef(routineName, "Found broadcast - uuid:%v ip:%v port:%v", uuid, ip, port)

			newDomainGuard := domain.NewDomainGuard(identity.DomainIdentity{
				Address: identity.Address{
					IP:   ip,
					Port: port,
				},
				UUID:     uuid,
				Services: make(map[string]identity.ServiceIdentity),
			})

			// Add the new member
			d.domains.Store(uuid, newDomainGuard)

		case <-ctx.Done():
			break Loop
		}
	}

	system.LogRoutinef(routineName, "Stopping routine")
}

func (d *Dominion) routineCheck(ctx context.Context, routineName string, wait time.Duration, check func(context.Context)) {
	system.LogRoutinef(routineName, "Starting routine")
	ticker := time.NewTicker(wait)

Loop:
	for {
		select {
		case <-ticker.C:
			check(ctx)
		case <-ctx.Done():
			break Loop
		}
	}
	system.LogRoutinef(routineName, "Stopping routine")
}

func (d *Dominion) checkDomains(ctx context.Context) {
	d.domains.Range(func(uuid string, domainGuard *domain.DomainGuard) bool {
		domainGuard.LatchRead(func(domain *domain.Domain) error {
			if time.Since(domain.LastContact) > config.GetDominionConfig().DomainCheck*10 {
				// its been a while, make sure they are still alive
				go d.rpcHeartbeat(ctx, domainGuard)
			}
			return nil
		})
		return true
	})
}

func (d *Dominion) checkServices(ctx context.Context) {
	requiredServices := config.GetServicesConfig().GetRequiredServices()
	dependencies := make(map[string]int)

	d.domains.Range(func(uuid string, domainGuard *domain.DomainGuard) bool {
		domainGuard.LatchRead(func(domain *domain.Domain) error {

			// check for requiredServices
			for serviceType, serviceDef := range requiredServices {
				if _, ok := domain.Services[serviceType]; ok {
					continue // Already hosting service
				}
				if domain.HasTraits(serviceDef.Traits) {
					go d.rpcStartService(ctx, domainGuard, serviceType)
				}
			}

			// find service dependencies
			for serviceType := range domain.DomainIdentity.Services {
				for _, dependency := range config.GetServicesConfig().Services[serviceType].Dependencies {
					dependencies[dependency]++
				}
			}
			return nil
		})
		return true
	})

	// Check that dependencies exist
	for dependency := range dependencies {

		if len(d.findService(dependency)) == 0 {
			canidates := d.findServiceCanidates(dependency)
			if len(canidates) == 0 {
				system.Errorf("No canidates availible for %v", dependency)
				continue
			}

			// TODO Handle multiples
			canidate := canidates[0]
			domainGuard, ok := d.domains.Load(canidate.UUID)
			if !ok {
				system.Errorf("Viable canidate no longer availible for %v", dependency)
				continue
			}

			go d.rpcStartService(ctx, domainGuard, dependency)
		}
	}
}

func (d *Dominion) hostDominion(ctx context.Context) error {
	routineName := "hostDominion"
	system.LogRoutinef(routineName, "Starting routine")

	address := fmt.Sprintf("%s:%v", "", d.Address.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("hostDominion() Failed to listen: %w", err)
	}

	server := grpc.NewServer()
	pb.RegisterDominionServer(server, d)
	// Register reflection service on gRPC server.
	go func() {
		<-ctx.Done()
		server.GracefulStop()
		system.LogRoutinef(routineName, "Stopped grpc server gracefully.")
	}()

	reflection.Register(server)
	err = server.Serve(lis)
	system.LogRoutinef(routineName, "Stopping routine")
	return err
}
