package service

import (
	"context"
	"errors"
	"math"

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
	serviceUUID := request.UUID
	foundPreviousDevice := false
	var err error
	var ctx context.Context
	l.Subscribers.Range(func(sub *Subscriber) bool {
		if sub.Device.GetID() == serviceUUID && sub.Device.GetType() == serviceType {
			if sub.IsConnected() {
				err = errors.New("UUID and Type matched an existing subscriber, but connection was still active")
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

		system.Logf("Reconnected old Device %s!", serviceUUID)
	} else {
		// No previous device found, build new subscriber
		var device device.Device
		switch serviceType {
		case "npBar":
			device = npdevice.NewBar(
				serviceUUID,
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
		system.Logf("Added new Device %s!", serviceUUID)
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

	var err error
	l.Subscribers.Range(func(sub *Subscriber) bool {
		device := sub.Device
		if device.GetID() != pbDevice.GetUUID() {
			return true
		}
		if device.GetType() != pbDevice.GetType() {
			err = errors.New("Found matching UUID, but type did not match")
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

	parentUUID := request.ParentUUID
	childUUID := request.ChildUUID

	var targetNode node.Node
	l.Subscribers.Range(func(sub *Subscriber) bool {
		nodes := sub.Device.GetNodes()
		for _, n := range nodes {
			if n.GetID() == childUUID {
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

	err := l.NodeTree.Insert(parentUUID, targetNode)

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}

// DeleteNode deletes a device from the NodeTree
// implements pb.LightOrchestratorServer
func (l *LightOrch) DeleteNode(ctx context.Context, request *pb.DeleteNodeRequest) (*pb.Empty, error) {
	rpcName := "DeleteNode"
	system.LogRPCf(rpcName, "Receiving request")

	parentUUID := request.ParentUUID
	childUUID := request.ChildUUID

	err := l.NodeTree.Delete(parentUUID, childUUID)

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}
