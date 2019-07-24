package lightsub

import (
	"context"
	"fmt"
	"time"

	"github.com/jmbarzee/domain/services"
	pb "github.com/jmbarzee/domain/services/lightorchastrator/grpc"
	ws2811 "github.com/jmbarzee/rpi_ws281x/golang/stub"
	"google.golang.org/grpc"
)

const (
	displayFPS                = 60
	displayRate time.Duration = time.Second / displayFPS

	gpioPin    = 18
	brightness = 255
)

type (
	LightSub struct {
		services.Service

		Size      int
		LightPlan LightPlan
	}
)

func NewLightSub(port int, size int, serviceName string, domainPort int) *LightSub {
	return &LightSub{
		Service: services.Service{
			ServiceName: serviceName,
			Port:        port,
			DomainPort:  domainPort,
		},

		Size:      size,
		LightPlan: NewLightPlan(),
	}
}

func (s *LightSub) Run() {
	go s.SubscribeLights()
	s.DisplayLights()
	// TODO @jmbarzee see if we should return an error
}

func (s *LightSub) DisplayLights() {
	defer ws2811.Fini()
	err := ws2811.Init(gpioPin, s.Size, brightness)
	if err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(displayRate)

	for {
		<-ticker.C
		var last *LightChange
		next := s.LightPlan.Peek()

		// if change has past search the next one
		for next.Time.Before(time.Now()) {
			last = &next
			next = s.LightPlan.Advance()
		}

		if last != nil {
			for i, wrgb := range last.Lights {
				ws2811.SetLed(i, wrgb)
			}
			ws2811.Render()

		}

	}

}

func (s *LightSub) SubscribeLights() {
	for {
		var err error
		ctx := context.TODO()
		addrs, err := s.Locate(ctx, "lightOrchastrator")
		if err != nil {
			// TODO @jmbarzee handle error
			continue
		}

		conn, err := grpc.DialContext(ctx, addrs[0], grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			// TODO @jmbarzee handle error
			continue
		}

		request := &pb.SubscribeLightsRequest{
			ServiceName: s.ServiceName,
		}

		client := pb.NewLightOrcharstratorClient(conn)
		subLightsClient, err := client.SubscribeLights(ctx, request)
		if err != nil {
			// TODO @jmbarzee handle error
			continue
		}

		for {
			reply, err := subLightsClient.Recv()
			if err != nil {
				// TODO @jmbarzee handle error
				break
			}

			lightChange := s.convertDLRtoLightChange(reply)
			s.LightPlan.Add(lightChange)
		}

		// TODO @jmbarzee consider sending a close message
	}
}

func (s LightSub) convertDLRtoLightChange(req *pb.SubscribeLightsReply) LightChange {
	change := LightChange{
		Time: time.Unix(0, req.GetDisplayTime()),
	}
	for i, color := range req.GetColors() {
		if i == s.Size {
			break
			fmt.Printf("More than %v colors sent", s.Size)
		}
		change.Lights[i] = uint32(color)
	}
	return change
}
