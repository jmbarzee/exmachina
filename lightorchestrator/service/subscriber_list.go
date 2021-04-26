package service

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/jmbarzee/services/lightorchestrator/grpc"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	device.Device
	server pb.LightOrchestrator_SubscribeLightsServer
	kill   context.CancelFunc
}

// IsConnected returns true if the server and kill functions are non-nil
func (s Subscriber) IsConnected() bool {
	return s.server != nil && s.kill != nil
}

// Connect adds a connection to an existing subscriber
// and returns a context to be waited on while the connection should be held open
func (s *Subscriber) Connect(server pb.LightOrchestrator_SubscribeLightsServer) context.Context {
	var ctx context.Context
	s.server = server
	ctx, s.kill = context.WithCancel(context.Background())
	return ctx
}

// Disconnect Ends a subscribers connection
// The connection is ended by cancling the context
// which is held in the origininal subscriber grpc
func (s *Subscriber) Disconnect() error {
	if !s.IsConnected() {
		return errors.New("Subscriber is not connected, cannot Disconnect")
	}
	s.kill()
	s.server = nil
	s.kill = nil
	return nil
}

// DispatchRender sends lights after a subscriber's device renders them based on t
func (s Subscriber) DispatchRender(t time.Time) error {
	lights := s.Render(t)

	colors := make([]uint32, len(lights))
	for i, light := range lights {
		rgb := light.GetColor().RGB()
		colors[i] = rgb.ToUInt32RGBW()
	}

	timestamp := timestamppb.New(t)

	reply := &pb.SubscribeLightsReply{
		DisplayTime: timestamp,
		Colors:      colors,
	}
	return s.server.Send(reply)
}

// SubscriberList thread-safe list of subscribers
type SubscriberList struct {
	// RWMutex gates changes to the list
	rwmutex *sync.RWMutex
	// subs is the list of subscriber
	subs []Subscriber
}

// Range ranges over a SubscriberList
func (l SubscriberList) Range(f func(sub *Subscriber) bool) {
	l.rwmutex.Lock()
	for i := 0; i < len(l.subs); i++ {
		if !f(&l.subs[i]) {
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
