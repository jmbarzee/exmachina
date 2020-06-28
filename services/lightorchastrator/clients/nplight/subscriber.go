package nplight

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmbarzee/domain/services"
	pb "github.com/jmbarzee/domain/services/lightorchastrator/grpc"
	ws2811 "github.com/jmbarzee/rpi_ws281x/golang/stub"
	"google.golang.org/grpc"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second / displayFPS

	gpioPin    = 18
	brightness = 255
)

type (
	Subscriber struct {
		services.Service

		Size      int
		LightPlan LightPlan
	}
)

func NewSubscriber(port int, domainPort int, logger *log.Logger, size int) *Subscriber {
	logger.Printf("Subscriber built!")
	return &Subscriber{
		Service: services.Service{
			ServiceName: "",
			Port:        port,
			DomainPort:  domainPort,
		},

		Size:      size,
		LightPlan: NewLightPlan(),
	}
}

func (s *Subscriber) Run() {
	s.Logger.Printf("Running Subscriber...")
	go s.SubscribeLights()
	s.DisplayLights()
}

func (s *Subscriber) DisplayLights() {
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

func (s *Subscriber) SubscribeLights() {
	for {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		addrs, err := s.Locate(ctx, "lightOrchastrator")
		if err != nil {
			s.Logger.Printf("Error locating lightOrchastrator: %v", err.Error())
			continue
		}

		conn, err := grpc.DialContext(ctx, addrs[0], grpc.WithInsecure(), grpc.WithBlock())
		defer conn.Close()
		if err != nil {
			s.Logger.Printf("Error dialing lightOrchastrator: %v", err.Error())
			continue
		}

		request := &pb.SubscribeLightsRequest{
			ServiceName: s.ServiceName,
		}

		client := pb.NewLightOrcharstratorClient(conn)
		subLightsClient, err := client.SubscribeLights(ctx, request)
		if err != nil {
			s.Logger.Printf("Error subscribing to lightOrchastrator: %v", err.Error())
			continue
		}

		for {
			reply, err := subLightsClient.Recv()
			if err != nil {
				s.Logger.Printf("Error receving reply from lightOrchastrator: %v", err.Error())
				break
			}

			lightChange := s.convertDLRtoLightChange(reply)
			s.LightPlan.Add(lightChange)
		}

		if err = conn.Close(); err != nil {
			s.Logger.Printf("Error closing connection to lightOrchastrator: %v", err.Error())
		}
	}
}

func (s Subscriber) convertDLRtoLightChange(req *pb.SubscribeLightsReply) LightChange {
	change := LightChange{
		Time: time.Unix(0, req.GetDisplayTime()),
	}
	for i, color := range req.GetColors() {
		if i == s.Size {
			break
			s.Logger.Printf("More than %v colors sent", s.Size)
		}
		change.Lights[i] = uint32(color)
	}
	return change
}
