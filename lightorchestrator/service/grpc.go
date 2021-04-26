package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jmbarzee/dominion/system"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/device/npdevice"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/pbconvert"
	"github.com/jmbarzee/space"
)

// SubscribeLights requests a stream of lights
// implements pb.LightOrchestratorServer
func (l *LightOrch) SubscribeLights(request *pb.SubscribeLightsRequest, server pb.LightOrchestrator_SubscribeLightsServer) error {
	rpcName := "SubscribeLights"
	system.LogRPCf(rpcName, "Receiving request")

	serviceType := request.Type
	serviceID, err := uuid.FromBytes(request.ID)
	if err != nil {
		system.Errorf("Failed to parse id %w", err)
		return err
	}

	foundPreviousDevice := false
	var ctx context.Context
	l.Subscribers.Range(func(sub *Subscriber) bool {
		if sub.Device.GetID() == serviceID && sub.Device.GetType() == serviceType {
			if sub.IsConnected() {
				err = errors.New("ID and Type matched an existing subscriber, but connection was still active")
				return false
			}
			foundPreviousDevice = true
			ctx = sub.Connect(server)
			return false
		}
		return true
	})

	if foundPreviousDevice {
		if err != nil {
			system.Errorf("Failed to reconnect device %w", err)
			return err
		}

		system.Logf("Reconnected old Device %s!", serviceID)
	} else {
		// No previous device found, build new subscriber
		var device device.Device
		switch serviceType {
		case "npBar":
			device = npdevice.NewBar(
				serviceID,
				space.Cartesian{X: 0, Y: 0, Z: 0},
				space.Spherical{R: 1, P: 0, T: 0},
				space.Spherical{R: 1, P: math.Pi / 2, T: 0},
			)
			// TODO @jmbarzee add other devices for start up here
		}
		if device == nil {
			return errors.New("Unrecognized Service Name")
		}
		sub := Subscriber{Device: device}
		ctx = (&sub).Connect(server)

		l.Subscribers.Append(sub)
		system.Logf("Added new Device %s!", serviceID)
	}

	// hold connection open until it is ended elsewhere
	<-ctx.Done()
	system.LogRPCf(rpcName, "Ending stream")
	return nil
}

// GetDevices returns the DeviceNode hierarchy and all subscribed devices
// implements pb.LightOrchestratorServer
func (l *LightOrch) GetDevices(ctx context.Context, request *pb.Empty) (*pb.GetDevicesReply, error) {
	rpcName := "GetDevices"
	system.LogRPCf(rpcName, "Receiving request")
	pbDeviceNodes := l.NodeTree.ToPBDeviceNode()
	pbDevices := make([]*pb.Device, 0)
	l.Subscribers.Range(func(sub *Subscriber) bool {
		pbDevices = append(pbDevices, pbconvert.NewPBDevice(sub.Device))
		return true
	})

	reply := &pb.GetDevicesReply{
		DeviceNodeTree: pbDeviceNodes,
		Devices:        pbDevices,
	}
	system.LogRPCf(rpcName, "Sending reply")
	return reply, nil
}

// MoveDevice changes a devices location and orientation
// implements pb.LightOrchestratorServer
func (l *LightOrch) MoveDevice(ctx context.Context, request *pb.MoveDeviceRequest) (*pb.Empty, error) {
	rpcName := "MoveDevice"
	system.LogRPCf(rpcName, "Receiving request")

	pbDevice := request.Device
	deviceID, err := uuid.FromBytes(pbDevice.GetID())
	if err != nil {
		err = fmt.Errorf("failed to move device: %w", err)
		system.LogRPCf(rpcName, err.Error())
		return nil, err
	}

	l.Subscribers.Range(func(sub *Subscriber) bool {
		device := sub.Device
		if device.GetID() != deviceID {
			return true
		}
		if device.GetType() != pbDevice.GetType() {
			err = errors.New("Found matching ID, but type did not match")
			return false
		}
		device.SetLocation(pbconvert.NewVector(pbDevice.GetLocation()))
		device.SetOrientation(pbconvert.NewOrientation(pbDevice.GetOrientation()))
		return false
	})

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}

// InsertNode inserts a device into the NodeTree
// implements pb.LightOrchestratorServer
func (l *LightOrch) InsertNode(ctx context.Context, request *pb.InsertNodeRequest) (*pb.Empty, error) {
	rpcName := "InsertNode"
	system.LogRPCf(rpcName, "Receiving request")

	parentID, err := uuid.FromBytes(request.GetParentID())
	if err != nil {
		err = fmt.Errorf("failed to insert node : %w", err)
		system.LogRPCf(rpcName, err.Error())
		return nil, err
	}

	childID, err := uuid.FromBytes(request.GetChildID())
	if err != nil {
		err = fmt.Errorf("failed to insert node : %w", err)
		system.LogRPCf(rpcName, err.Error())
		return nil, err
	}

	var targetNode node.Node
	l.Subscribers.Range(func(sub *Subscriber) bool {
		nodes := sub.Device.GetNodes()
		for _, n := range nodes {
			if n.GetID() == childID {
				targetNode = n
			}
		}

		if targetNode == nil {
			return true
		}
		return false
	})

	if targetNode == nil {
		return nil, errors.New("Could not find specified Child")
	}

	err = l.NodeTree.Insert(parentID, targetNode)

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}

// DeleteNode deletes a device from the NodeTree
// implements pb.LightOrchestratorServer
func (l *LightOrch) DeleteNode(ctx context.Context, request *pb.DeleteNodeRequest) (*pb.Empty, error) {
	rpcName := "DeleteNode"
	system.LogRPCf(rpcName, "Receiving request")

	parentID, err := uuid.FromBytes(request.GetParentID())
	if err != nil {
		err = fmt.Errorf("failed to insert node : %w", err)
		system.LogRPCf(rpcName, err.Error())
		return nil, err
	}

	childID, err := uuid.FromBytes(request.GetChildID())
	if err != nil {
		err = fmt.Errorf("failed to insert node : %w", err)
		system.LogRPCf(rpcName, err.Error())
		return nil, err
	}

	err = l.NodeTree.Delete(parentID, childID)

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}
