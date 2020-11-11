package npdevice

import (
	"math"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/node/npnode"
	"github.com/jmbarzee/space"
)

// ChandelierSmall is a Small Chandelier (2 rings)
type ChandelierSmall struct {
	device.Basic

	*space.Object

	SmallRing *npnode.Ring
	LargeRing *npnode.Ring
}

var _ device.Device = (*ChandelierSmall)(nil)

// NewChandelierSmall returns a new Small Chandelier
func NewChandelierSmall(uuid string, top space.Cartesian, theta float64) ChandelierSmall {
	orientation := space.NewSpherical(1, theta, 0)
	rotation := space.NewSpherical(1, theta, math.Pi/2)
	d := ChandelierSmall{
		Basic:  device.NewBasic(uuid),
		Object: space.NewObject(top, orientation, rotation),
	}

	l, o, r := d.getSpaceDataForSmallRing()
	d.SmallRing = npnode.NewRing(smallRingRadius, l, o, r)

	l, o, r = d.getSpaceDataForLargeRing()
	d.LargeRing = npnode.NewRing(largeRingRadius, l, o, r)

	return d
}

// GetNodes returns all the Nodes which the device holds
func (d ChandelierSmall) GetNodes() []node.Node {
	return []node.Node{
		d.SmallRing,
		d.LargeRing,
	}
}

// SetLocation changes the physical location of the device
func (d ChandelierSmall) SetLocation(v space.Cartesian) {
	d.Object.SetLocation(v)
	d.updateRingLocations()
}

// SetOrientation changes the physical orientation of the device
// o.Phi will be forced to 0 (because it hangs strait down)
func (d ChandelierSmall) SetOrientation(o space.Spherical) {
	o.P = 0
	d.Object.SetOrientation(o)
	d.updateRingLocations()
}

// SetRotation changes the physical rotation of the device
func (d ChandelierSmall) SetRotation(o space.Spherical) {
	d.Object.SetRotation(o)
	d.updateRingLocations()
}

func (d ChandelierSmall) updateRingLocations() {
	l, o, r := d.getSpaceDataForSmallRing()
	d.SmallRing.SetLocation(l)
	d.SmallRing.SetOrientation(o)
	d.SmallRing.SetRotation(r)

	l, o, r = d.getSpaceDataForLargeRing()
	d.LargeRing.SetLocation(l)
	d.LargeRing.SetOrientation(o)
	d.LargeRing.SetRotation(r)
}

func (d ChandelierSmall) getSpaceDataForSmallRing() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: smallRingHeight1}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(smallRingTilt1).Rotate(smallRingRotation1)
	rotation := d.GetRotation().Tilt(smallRingTilt1).Rotate(smallRingRotation1)

	return location, orientation, rotation
}

func (d ChandelierSmall) getSpaceDataForLargeRing() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: largeRingHeight1}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(largeRingTilt1).Rotate(largeRingRotation1)
	rotation := d.GetRotation().Tilt(largeRingTilt1).Rotate(largeRingRotation1)

	return location, orientation, rotation
}

// Render calls render on each of the rings and then appends all the lights
func (d ChandelierSmall) Render(t time.Time) []ifaces.Light {
	allLights := append(d.SmallRing.Render(t), d.LargeRing.Render(t)...)
	return allLights
}

// GetType returns the type
func (d ChandelierSmall) GetType() string {
	return "npChandelierSmall"
}
