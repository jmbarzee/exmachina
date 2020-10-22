package nplight

import (
	"context"
	"time"
)

func (l *NPLight) updateLights(ctx context.Context, t time.Time) {
	// Advance the light plan
	next := l.LightPlan.Advance(t)
	if next != nil {
		for i, wrgb := range next.Lights {
			l.Strip.Leds(0)[i] = wrgb
		}
		l.Strip.Render()
	}
}
