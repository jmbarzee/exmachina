package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect"
	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/space"
)

// ChandelierSmall is a Small Chandelier (2 rings)
type ChandelierSmall struct {
	SmallRing *Ring
	LargeRing *Ring
	Top       space.Vector // Mounting location for the chandilier
}

// NewChandelierSmall returns a new Small Chandelier
func NewChandelierSmall(top space.Vector, theta float64) ChandelierSmall {
	center := space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - .66,
	}
	smallRing := NewRing(center, 0.7, 0.0+theta, math.Pi/6)
	largeRing := NewRing(center, 1.3, math.Pi/2+theta, math.Pi/6)
	return ChandelierSmall{
		SmallRing: smallRing,
		LargeRing: largeRing,
		Top:       top,
	}
}

// Allocate takes Vibes and Distributes them to the rings
func (c ChandelierSmall) Allocate(feeling effect.Vibe) {
	newVibe := feeling.Stabalize()
	c.SmallRing.Allocate(newVibe)
	c.LargeRing.Allocate(newVibe)
}

// Render calls render on each of the rings and then appends all the lights
func (c ChandelierSmall) Render(t time.Time) []*effect.Light {
	allLights := append(c.SmallRing.Render(t), c.LargeRing.Render(t)...)
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (c ChandelierSmall) PruneEffects(t time.Time) {
	c.SmallRing.PruneEffects(t)
	c.LargeRing.PruneEffects(t)
}
