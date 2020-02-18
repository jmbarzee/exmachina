package lightorchastrator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jmbarzee/domain/services"
	"github.com/jmbarzee/domain/services/lightorchastrator/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/effect/devices"
	pb "github.com/jmbarzee/domain/services/lightorchastrator/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second / displayFPS
)

type (
	LightOrch struct {
		services.Service
		sync.Mutex
		Subs []Subscription
	}

	Subscription struct {
		Body   LightSub
		Server pb.LightOrcharstrator_SubscribeLightsServer
		Kill   context.CancelFunc
	}
)

func NewLightOrch(port int, domainPort int) *LightOrch {
	return &LightOrch{
		Service: services.Service{
			ServiceName: "lightOrchastrator",
			Port:        port,
			DomainPort:  domainPort,
		},
		Subs: []Subscription{},
	}
}

func (l *LightOrch) Run() {
	go l.listen()
	l.orchastrate()

}

func (l *LightOrch) orchastrate() {
	ticker := time.NewTicker(displayRate)

	for <-ticker {
		for _, sub := l.Subs {
			
		}
	}

}

func (l *LightOrch) listen() {
	address := fmt.Sprintf("%s:%v", "", l.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterLightOrcharstratorServer(server, l)

	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

}

func (l *LightOrch) SubscribeLights(request *pb.SubscribeLightsRequest, server pb.LightOrcharstrator_SubscribeLightsServer) error {
	serviceName := request.ServiceName
	var body LightSub
	switch serviceName {
	case "neoPixelBar":
		body = devices.NewNeoPixelBar()
		// TODO @jmbarzee add other devices for start up here
	}
	if body == nil {
		return errors.New("Unrecognized Service Name")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	sub := Subscription{
		Body:   body,
		Server: server,
		Kill:   cancelFunc,
	}
	l.Lock()
	l.Subs = append(l.Subs, sub)
	l.Unlock()

	// hold connection open until it is ended elsewhere 
	<-ctx.Done()
	return nil //TODO @jmbarzee consider sending error message
}
