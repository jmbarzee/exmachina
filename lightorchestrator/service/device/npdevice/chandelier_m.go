package npdevice

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/node/npnode"
	"github.com/jmbarzee/space"
)

const (
	smallRingRadius = 0.7

	smallRingHeight1   = -0.6
	smallRingTilt1     = 2 / 6 * math.Pi
	smallRingRotation1 = 0 / 6 * math.Pi

	smallRingHeight2   = -1.2
	smallRingTilt2     = 2 / 6 * math.Pi
	smallRingRotation2 = 3 / 6 * math.Pi

	largeRingRadius = 1.3

	largeRingHeight1   = -0.6
	largeRingTilt1     = 2 / 6 * math.Pi
	largeRingRotation1 = 6 / 6 * math.Pi

	largeRingHeight2   = -1.2
	largeRingTilt2     = 2 / 6 * math.Pi
	largeRingRotation2 = 9 / 6 * math.Pi
)

// ChandelierMedium is a Medium Chandelier (4 rings)
type ChandelierMedium struct {
	device.Basic

	*space.Object

	SmallRing1 *npnode.Ring
	SmallRing2 *npnode.Ring
	LargeRing1 *npnode.Ring
	LargeRing2 *npnode.Ring
}

var _ device.Device = (*ChandelierMedium)(nil)

// NewChandelierMedium returns a new Medium Chandelier
func NewChandelierMedium(id uuid.UUID, top space.Cartesian, theta float64) ChandelierMedium {
	orientation := space.NewSpherical(1, theta, 0)
	rotation := space.NewSpherical(1, theta, math.Pi/2)
	d := ChandelierMedium{
		Basic:  device.NewBasic(id),
		Object: space.NewObject(top, orientation, rotation),
	}

	l, o, r := d.getSpaceDataForSmallRing1()
	d.SmallRing1 = npnode.NewRing(smallRingRadius, l, o, r)

	l, o, r = d.getSpaceDataForSmallRing2()
	d.SmallRing2 = npnode.NewRing(smallRingRadius, l, o, r)

	l, o, r = d.getSpaceDataForLargeRing1()
	d.LargeRing1 = npnode.NewRing(largeRingRadius, l, o, r)

	l, o, r = d.getSpaceDataForLargeRing2()
	d.LargeRing2 = npnode.NewRing(largeRingRadius, l, o, r)

	return d
}

// GetNodes returns all the Nodes which the device holds
func (d ChandelierMedium) GetNodes() []node.Node {
	return []node.Node{
		d.SmallRing1,
		d.SmallRing2,
		d.LargeRing1,
		d.LargeRing2,
	}
}

// SetLocation changes the physical location of the device
func (d ChandelierMedium) SetLocation(v space.Cartesian) {
	d.Object.SetLocation(v)
	d.updateRingLocations()
}

// SetOrientation changes the physical orientation of the device
// o.Phi will be forced to 0 (because it hangs strait down)
func (d ChandelierMedium) SetOrientation(o space.Spherical) {
	o.P = 0
	d.Object.SetOrientation(o)
	d.updateRingLocations()
}

// SetRotation changes the physical rotation of the device
func (d ChandelierMedium) SetRotation(o space.Spherical) {
	d.Object.SetRotation(o)
	d.updateRingLocations()
}

func (d ChandelierMedium) updateRingLocations() {
	l, o, r := d.getSpaceDataForSmallRing1()
	d.SmallRing1.SetLocation(l)
	d.SmallRing1.SetOrientation(o)
	d.SmallRing1.SetRotation(r)

	l, o, r = d.getSpaceDataForSmallRing2()
	d.SmallRing2.SetLocation(l)
	d.SmallRing2.SetOrientation(o)
	d.SmallRing2.SetRotation(r)

	l, o, r = d.getSpaceDataForLargeRing1()
	d.LargeRing1.SetLocation(l)
	d.LargeRing1.SetOrientation(o)
	d.LargeRing1.SetRotation(r)

	l, o, r = d.getSpaceDataForLargeRing2()
	d.LargeRing2.SetLocation(l)
	d.LargeRing2.SetOrientation(o)
	d.LargeRing2.SetRotation(r)
}

func (d ChandelierMedium) getSpaceDataForSmallRing1() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: smallRingHeight1}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(smallRingTilt1).Rotate(smallRingRotation1)
	rotation := d.GetRotation().Tilt(smallRingTilt1).Rotate(smallRingRotation1)

	return location, orientation, rotation
}

func (d ChandelierMedium) getSpaceDataForSmallRing2() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: smallRingHeight2}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(smallRingTilt2).Rotate(smallRingRotation2)
	rotation := d.GetRotation().Tilt(smallRingTilt2).Rotate(smallRingRotation2)

	return location, orientation, rotation
}

func (d ChandelierMedium) getSpaceDataForLargeRing1() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: largeRingHeight1}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(largeRingTilt1).Rotate(largeRingRotation1)
	rotation := d.GetRotation().Tilt(largeRingTilt1).Rotate(largeRingRotation1)

	return location, orientation, rotation
}

func (d ChandelierMedium) getSpaceDataForLargeRing2() (space.Cartesian, space.Spherical, space.Spherical) {
	locationTransformation := space.Cartesian{X: 0, Y: 0, Z: largeRingHeight2}.TranslationMatrix()
	location := d.GetLocation().Transform(locationTransformation).Cartesian()

	orientation := d.GetOrientation().Tilt(largeRingTilt2).Rotate(largeRingRotation2)
	rotation := d.GetRotation().Tilt(largeRingTilt2).Rotate(largeRingRotation2)

	return location, orientation, rotation
}

// Render calls render on each of the rings and then appends all the lights
func (d ChandelierMedium) Render(t time.Time) []ifaces.Light {
	allLights := []ifaces.Light{}
	allLights = append(allLights, d.SmallRing1.Render(t)...)
	allLights = append(allLights, d.SmallRing2.Render(t)...)
	allLights = append(allLights, d.LargeRing1.Render(t)...)
	allLights = append(allLights, d.LargeRing2.Render(t)...)
	return allLights
}

// GetType returns the type
func (d ChandelierMedium) GetType() string {
	return "npChandelierMedium"
}
