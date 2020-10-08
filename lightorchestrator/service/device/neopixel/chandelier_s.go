package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/shared"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe"
)

// ChandelierSmall is a Small Chandelier (2 rings)
type ChandelierSmall struct {
	device.BasicDevice
	SmallRing *Ring
	LargeRing *Ring
	Top       space.Vector // Mounting location for the chandilier
}

// NewChandelierSmall returns a new Small Chandelier
func NewChandelierSmall(top space.Vector, theta float32) ChandelierSmall {
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
func (c ChandelierSmall) Allocate(feeling vibe.Vibe) {
	newVibe := feeling.Stabilize()
	c.SmallRing.Allocate(newVibe)
	c.LargeRing.Allocate(newVibe)
}

// Render calls render on each of the rings and then appends all the lights
func (c ChandelierSmall) Render(t time.Time) []shared.Light {
	allLights := append(c.SmallRing.Render(t), c.LargeRing.Render(t)...)
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (c ChandelierSmall) PruneEffects(t time.Time) {
	c.SmallRing.PruneEffects(t)
	c.LargeRing.PruneEffects(t)
}
