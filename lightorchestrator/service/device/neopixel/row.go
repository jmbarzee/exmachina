package neopixel

import (
	"sort"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
)

const (
	ledsPerMeter = 60

	// density is assumed to be 60 leds per meter
	distPerLED = 1.0 / float64(ledsPerMeter)
)

// Row holds basic functionality for any neoPixelDevice
// Row implements effect.Allocater
// Row partially implements effect.Device with PruneEffect
type Row struct {
	// Length is the number of LEDs in the Row
	Length int

	// Effects is the array of effects from materializing Vibes in calls to Allocate
	Effects []ifaces.Effect

	GetLights func() []light.Light
}

// NewRow creates a new Row
func NewRow(
	length int,
	getLights func() []light.Light,
) *Row {
	return &Row{
		Length:    length,
		Effects:   make([]ifaces.Effect, 0),
		GetLights: getLights,
	}
}

// PruneEffects removes all effects which have ended before a time t
func (d *Row) PruneEffects(t time.Time) {
	unEndedEffects := make([]ifaces.Effect, 0, len(d.Effects))
	for _, e := range d.Effects {
		if e.End().Before(t) {
			continue
		}

		unEndedEffects = append(unEndedEffects)
	}

	d.Effects = unEndedEffects
}

// Render uses the stored effects from allocate(feeling)s to produce an array of lights
func (d *Row) Render(t time.Time) []light.Light {
	runningEffects := make([]ifaces.Effect, 0)
	for _, f := range d.Effects {
		if f.Start().After(time.Now()) {
			// since the list is ordered we can assume there are no more running effects
			break
		} else if f.End().After(time.Now()) {
			runningEffects = append(runningEffects, f)
		}
	}

	sort.Sort(byPriority(runningEffects))

	lights := d.GetLights()

	for _, e := range runningEffects {
		lights = e.Render(t, lights)
	}

	return lights
}

// GetChildren returns all groups under the GroupOption
func (d Row) GetChildren() []device.DeviceNode {
	return []device.DeviceNode{}
}

// Allocate takes Vibes and Materializes them into effects
// This is the bottom of the Allocater hiarchy for neoPixels
func (d *Row) Allocate(feeling ifaces.Vibe) {
	newEffects := feeling.Materialize()
	d.Effects = append(d.Effects, newEffects...)
	sort.Sort(byStartTime(d.Effects))
}

type byStartTime []ifaces.Effect

func (st byStartTime) Len() int           { return len(st) }
func (st byStartTime) Swap(i, j int)      { st[i], st[j] = st[j], st[i] }
func (st byStartTime) Less(i, j int) bool { return st[i].Start().Before(st[j].Start()) }

type byPriority []ifaces.Effect

func (p byPriority) Len() int           { return len(p) }
func (p byPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byPriority) Less(i, j int) bool { return p[i].Priotity() < p[j].Priotity() }
