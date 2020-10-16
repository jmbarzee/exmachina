package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// ChandelierSmall is a Small Chandelier (2 rings)
type ChandelierSmall struct {
	device.BasicDevice

	// Top is the mounting location for the chandilier
	Top space.Vector

	SmallRing *Ring
	LargeRing *Ring
}

// NewChandelierSmall returns a new Small Chandelier
func NewChandelierSmall(top space.Vector, theta float32) ChandelierSmall {

	center := space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - .6,
	}
	orientation := space.Orientation{
		Theta: theta + math.Pi/6,
	}
	smallRing := NewRing(smallRingRadius, center, orientation)

	orientation = orientation.Rotate(math.Pi / 2)
	largeRing := NewRing(largeRingRadius, center, orientation)

	return ChandelierSmall{
		SmallRing: smallRing,
		LargeRing: largeRing,
		Top:       top,
	}
}

// Allocate takes Vibes and Distributes them to the rings
func (d ChandelierSmall) Allocate(feeling ifaces.Vibe) {
	newVibe := feeling.Stabilize()
	d.SmallRing.Allocate(newVibe)
	d.LargeRing.Allocate(newVibe)
}

// Render calls render on each of the rings and then appends all the lights
func (d ChandelierSmall) Render(t time.Time) []light.Light {
	allLights := append(d.SmallRing.Render(t), d.LargeRing.Render(t)...)
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (d ChandelierSmall) PruneEffects(t time.Time) {
	d.SmallRing.PruneEffects(t)
	d.LargeRing.PruneEffects(t)
}

// GetType returns the type
func (d ChandelierSmall) GetType() string {
	return "npChandelierSmall"
}
