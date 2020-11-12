package service

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/span"
)

const (
	tickLength = time.Second * 5
)

func (l *LightOrch) dispatchRender(ctx context.Context, t time.Time) {
	l.Subscribers.Range(func(sub *Subscriber) bool {
		if !sub.IsConnected() {
			return true
		}

		if err := sub.DispatchRender(t); err != nil {
			system.Errorf("Failed to dispatch Render: %w", err)
			sub.Disconnect()
		}
		return true
	})
}

func (l *LightOrch) allocateVibe(ctx context.Context, t time.Time) {
	v := &vibe.Basic{
		Span: span.Span{
			StartTime: t.Add(tickLength),
			EndTime:   t.Add(tickLength * 2),
		},
	}
	l.NodeTree.Allocate(v)
}
