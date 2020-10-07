package nplight

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmbarzee/dominion/identity"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/services/lightorchestrator/clients/nplight/lightplan"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/dominion/system"
	"google.golang.org/grpc"
)

func (l *NPLight) rpcSubscribeLights(ctx context.Context, lightOrchestrator identity.ServiceIdentity) error {
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
		Type: l.Type,
		UUID: l.UUID,
	}

	system.LogRPCf(rpcName, "Sending request")
	client := pb.NewLightOrchestratorClient(conn)
	subLightsClient, err := client.SubscribeLights(ctx, request)
	if err != nil {
		return fmt.Errorf("Error subscribing to lightOrchestrator: %w", err)
	}
	system.LogRPCf(rpcName, "Recieved stream client")

	for {
		reply, err := subLightsClient.Recv()
		if err != nil {
			if err = conn.Close(); err != nil {
				return fmt.Errorf("Error closing connection to lightOrchestrator: %v", err)
			}
			return fmt.Errorf("Error receving reply from lightOrchestrator: %v", err)
		}

		lightChange, err := l.convertDLRtoLightChange(reply)
		if err != nil {
			system.Errorf("Could not convert to LightChange: %w", err)
			continue
		}
		l.LightPlan.Add(lightChange)
	}
}

func (l *NPLight) convertDLRtoLightChange(reply *pb.SubscribeLightsReply) (lightplan.LightChange, error) {
	t, err := ptypes.Timestamp(reply.GetDisplayTime())
	if err != nil {
		return lightplan.LightChange{}, err
	}

	change := lightplan.LightChange{
		Time:   t,
		Lights: make([]uint32, l.Size),
	}
	for i, color := range reply.GetColors() {
		if i == l.Size {
			return lightplan.LightChange{}, fmt.Errorf("Expected %v colors, got %v", l.Size, len(reply.GetColors()))
		}
		change.Lights[i] = uint32(color)
	}
	return change, nil
}
