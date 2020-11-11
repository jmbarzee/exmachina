package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
)

// NewStructs produces a SubscriberList and NodeTree which share an underlying lock for thread safety
func NewStructs() (*SubscriberList, *NodeTree) {
	root := node.NewGroupOption()
	rwmutex := sync.RWMutex{}

	return &SubscriberList{
			rwmutex: &rwmutex,
			subs:    []Subscriber{},
		},
		&NodeTree{
			root:    root,
			rwmutex: &rwmutex,
		}
}

// Subscriber represents a light service which has subscribed to light updates
// from the LightOrchestrator
type Subscriber struct {
	Device device.Device
	Server pb.LightOrchestrator_SubscribeLightsServer
	Kill   context.CancelFunc
}

// DispatchRender sends lights after a subscriber's device renders them based on t
func (s Subscriber) DispatchRender(t time.Time) error {
	lights := s.Device.Render(t)

	colors := make([]uint32, len(lights))
	for i, light := range lights {
		colors[i] = light.GetColor().ToRGBA().ToUInt32WGRB()
	}

	timestamp, err := ptypes.TimestampProto(t)
	if err != nil {
		return fmt.Errorf("Failed to create timestamp: %w", err)
	}

	reply := &pb.SubscribeLightsReply{
		DisplayTime: timestamp,
		Colors:      colors,
	}
	return s.Server.Send(reply)
}

// SubscriberList thread-safe list of subscribers
type SubscriberList struct {
	// RWMutex gates changes to the list
	rwmutex *sync.RWMutex
	// subs is the list of subscriber
	subs []Subscriber
}

// Range ranges over a SubscriberList
func (l SubscriberList) Range(f func(sub Subscriber) bool) {
	l.rwmutex.Lock()
	for _, sub := range l.subs {
		if !f(sub) {
			break
		}
	}
	l.rwmutex.Unlock()
}

// Append appends a subscriber to a SubscriberList
func (l *SubscriberList) Append(sub Subscriber) {
	l.rwmutex.Lock()
	l.subs = append(l.subs, sub)
	l.rwmutex.Unlock()
}
