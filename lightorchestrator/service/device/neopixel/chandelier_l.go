package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

const (
	smallRingRadius = 0.7
	largeRingRadius = 1.3
)

// ChandelierLarge is a Large Chandelier (6 rings)
type ChandelierLarge struct {
	device.BasicDevice

	// Top is the mounting location for the chandilier
	Top space.Vector

	SmallRings []*Ring
	LargeRings []*Ring

	Groupings device.AllocGroupOption
}

// NewChandelierLarge returns a new Large Chandelier
func NewChandelierLarge(top space.Vector, theta float32) ChandelierLarge {
	smallRings := make([]*Ring, 3)
	largeRings := make([]*Ring, 3)

	center := space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - .6,
	}
	orientation := space.Orientation{
		Theta: theta + math.Pi/6,
	}
	smallRings[0] = NewRing(smallRingRadius, center, orientation)

	orientation = orientation.Rotate(math.Pi / 2)
	largeRings[0] = NewRing(largeRingRadius, center, orientation)

	center = space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - 1.0,
	}
	smallRings[1] = NewRing(smallRingRadius, center, orientation)

	orientation = orientation.Rotate(math.Pi / 2)
	largeRings[1] = NewRing(largeRingRadius, center, orientation)

	center = space.Vector{
		X: top.X,
		Y: top.Y,
		Z: top.Z - 1.4,
	}
	smallRings[2] = NewRing(smallRingRadius, center, orientation)

	orientation = orientation.Rotate(math.Pi / 2)
	largeRings[2] = NewRing(largeRingRadius, center, orientation)

	groupings := device.NewAllocGroupOption(
		device.NewAllocGroup(
			device.NewAllocGroup(
				smallRings[0],
				smallRings[1],
				smallRings[2],
			),
			device.NewAllocGroup(
				largeRings[0],
				largeRings[1],
				largeRings[2],
			),
		),
		device.NewAllocGroup(
			device.NewAllocGroup(
				smallRings[0],
				largeRings[0],
			),
			device.NewAllocGroup(
				smallRings[1],
				largeRings[1],
			),
			device.NewAllocGroup(
				smallRings[2],
				largeRings[2],
			),
		),
	)

	return ChandelierLarge{
		SmallRings: smallRings,
		LargeRings: largeRings,
		Top:        top,
		Groupings:  groupings,
	}
}

// Allocate takes Vibes and Distributes them to the rings
func (d ChandelierLarge) Allocate(vibe ifaces.Vibe) {
	d.Groupings.Allocate(vibe)
}

// Render calls render on each of the rings and then appends all the lights
func (d ChandelierLarge) Render(t time.Time) []light.Light {
	allLights := []light.Light{}
	for i := 0; i < 2; i++ {
		smallLights := d.SmallRings[i].Render(t)
		allLights = append(allLights, smallLights...)

		largeLights := d.LargeRings[i].Render(t)
		allLights = append(allLights, largeLights...)
	}
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (d ChandelierLarge) PruneEffects(t time.Time) {
	for _, smallRing := range d.SmallRings {
		smallRing.PruneEffects(t)
	}
	for _, largeRing := range d.LargeRings {
		largeRing.PruneEffects(t)
	}
}

// GetType returns the type
func (d ChandelierLarge) GetType() string {
	return "npChandelierLarge"
}
