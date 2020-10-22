package neopixel

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

// ChandelierMedium is a Medium Chandelier (4 rings)
type ChandelierMedium struct {
	device.BasicDevice

	// Top is the mounting location for the chandilier
	Top space.Vector

	SmallRings []*Ring
	LargeRings []*Ring

	Groupings device.AllocGroupOption
}

// NewChandelierMedium returns a new Medium Chandelier
func NewChandelierMedium(top space.Vector, theta float64) ChandelierMedium {
	smallRings := make([]*Ring, 2)
	largeRings := make([]*Ring, 2)

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

	groupings := device.NewAllocGroupOption(
		device.NewAllocGroup(
			device.NewAllocGroup(
				smallRings[0],
				smallRings[1],
			),
			device.NewAllocGroup(
				largeRings[0],
				largeRings[1],
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
		),
	)

	return ChandelierMedium{
		SmallRings: smallRings,
		LargeRings: largeRings,
		Top:        top,
		Groupings:  groupings,
	}
}

// Allocate takes Vibes and Distributes them to the rings
func (d ChandelierMedium) Allocate(vibe ifaces.Vibe) {
	d.Groupings.Allocate(vibe)
}

// Render calls render on each of the rings and then appends all the lights
func (d ChandelierMedium) Render(t time.Time) []light.Light {
	allLights := []light.Light{}
	for i := 0; i < 3; i++ {
		smallLights := d.SmallRings[i].Render(t)
		allLights = append(allLights, smallLights...)

		largeLights := d.LargeRings[i].Render(t)
		allLights = append(allLights, largeLights...)
	}
	return allLights
}

// PruneEffects removes all effects from the rigns which have ended before a time t
func (d ChandelierMedium) PruneEffects(t time.Time) {
	for _, smallRing := range d.SmallRings {
		smallRing.PruneEffects(t)
	}
	for _, largeRing := range d.LargeRings {
		largeRing.PruneEffects(t)
	}
}

// GetType returns the type
func (d ChandelierMedium) GetType() string {
	return "npChandelierMedium"
}
