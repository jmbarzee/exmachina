package npsub

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/clients/npsub/lightplan"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"google.golang.org/grpc"
)

func (s *NPSub) rpcSubscribeLights(ctx context.Context, lightOrchestrator identity.ServiceIdentity) error {
	rpcName := "SubscribeLights"
	conn, err := grpc.DialContext(
		ctx,
		lightOrchestrator.Address.String(),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(config.DefaultServiceDialTimeout))
	if err != nil {
		return fmt.Errorf("Error dialing lightOrchestrator: %w", err)
	}
	defer conn.Close()

	request := &pb.SubscribeLightsRequest{
		Type: s.Type,
		UUID: s.UUID,
	}

	system.LogRPCf(rpcName, "Sending request")
	client := pb.NewLightOrchestratorClient(conn)
	subLightsClient, err := client.SubscribeLights(ctx, request)
	if err != nil {
		return fmt.Errorf("Error subscribing to lightOrchestrator: %w", err)
	}
	system.LogRPCf(rpcName, "Received stream client")

	for {
		reply, err := subLightsClient.Recv()
		if err != nil {
			if err = conn.Close(); err != nil {
				return fmt.Errorf("Error closing connection to lightOrchestrator: %v", err)
			}
			return fmt.Errorf("Error receiving reply from lightOrchestrator: %v", err)
		}

		lightChange, err := s.convertDLRtoLightChange(reply)
		if err != nil {
			system.Errorf("Could not convert to LightChange: %w", err)
			continue
		}
		s.LightPlan.Add(lightChange)
	}
}

func (s *NPSub) convertDLRtoLightChange(reply *pb.SubscribeLightsReply) (lightplan.LightChange, error) {
	t, err := ptypes.Timestamp(reply.GetDisplayTime())
	if err != nil {
		return lightplan.LightChange{}, err
	}

	change := lightplan.LightChange{
		Time:   t,
		Lights: make([]uint32, s.Size),
	}
	for i, color := range reply.GetColors() {
		if i == s.Size {
			return lightplan.LightChange{}, fmt.Errorf("Expected %v colors, got %v", s.Size, len(reply.GetColors()))
		}
		change.Lights[i] = uint32(color)
	}
	return change, nil
}
