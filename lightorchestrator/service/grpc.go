package service

import (
	"context"
	"errors"

	"github.com/jmbarzee/dominion/system"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/device/neopixel"
	"github.com/jmbarzee/services/lightorchestrator/service/pbconvert"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

// SubscribeLights requests a stream of lights
// implements pb.LightOrchestratorServer
func (l *LightOrch) SubscribeLights(request *pb.SubscribeLightsRequest, server pb.LightOrchestrator_SubscribeLightsServer) error {
	rpcName := "SubscribeLights"
	system.LogRPCf(rpcName, "Receiving request")
	serviceType := request.Type
	serviceUUID := request.UUID
	var device device.Device
	switch serviceType {
	case "npBar":
		device = neopixel.NewBar(
			serviceUUID,
			space.Vector{X: 0, Y: 0, Z: 0},
			space.Orientation{Phi: 0, Theta: 0},
			space.Orientation{Phi: 0, Theta: 0},
		)
		// TODO @jmbarzee add other devices for start up here
	}
	if device == nil {
		return errors.New("Unrecognized Service Name")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	sub := Subscriber{
		Device: device,
		Server: server,
		Kill:   cancelFunc,
	}
	l.Subscribers.Append(sub)
	system.Logf("Added new Device %s!", serviceUUID)

	// hold connection open until it is ended elsewhere
	<-ctx.Done()
	system.LogRPCf(rpcName, "Ending stream")
	return nil //TODO @jmbarzee consider sending error message
}

// GetDevices returns the DeviceNode hierarchy and all subscribed devices
// implements pb.LightOrchestratorServer
func (l *LightOrch) GetDevices(ctx context.Context, request *pb.Empty) (*pb.GetDevicesReply, error) {
	rpcName := "GetDevices"
	system.LogRPCf(rpcName, "Receiving request")
	pbDeviceNodes := l.DeviceHierarchy.ToPBDeviceNode()
	pbDevices := make([]*pb.Device, 0)
	l.Subscribers.Range(func(sub Subscriber) bool {
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
	l.Subscribers.Range(func(sub Subscriber) bool {
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

// InsertDeviceInHierarchy inserts a device into the DeviceNode hierarchy
// implements pb.LightOrchestratorServer
func (l *LightOrch) InsertDeviceInHierarchy(ctx context.Context, request *pb.InsertDeviceInHierarchyRequest) (*pb.Empty, error) {
	rpcName := "InsertDeviceInHierarchy"
	system.LogRPCf(rpcName, "Receiving request")

	parentUUID := request.ParentUUID
	childUUID := request.ChildUUID

	var targetDevice device.Device
	l.Subscribers.Range(func(sub Subscriber) bool {
		device := sub.Device
		if device.GetID() != childUUID {
			return true
		}
		targetDevice = sub.Device
		return false
	})

	if targetDevice == nil {
		return nil, errors.New("Could not find specified Child")
	}

	err := l.DeviceHierarchy.Insert(parentUUID, targetDevice)

	system.LogRPCf(rpcName, "Sending reply")
	return &pb.Empty{}, err
}
