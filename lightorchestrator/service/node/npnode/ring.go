package npnode

import (
	"math"

	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/space"
)

// Ring is a representation of a ring of neopixels.
// Ring implements effect.Device
type Ring struct {
	// Row provides the implementation of effect.Allocater
	*Row

	node.Basic

	// Radius is the distance from the Center to any LED
	Radius float64

	*space.Object
}

var _ node.Node = (*Ring)(nil)

// NewRing creates a new Ring
func NewRing(
	radius float64,
	center space.Cartesian,
	orientation space.Spherical, // Orientation of center of ring, right hand rule
	rotation space.Spherical, // Orientation of First LED, forced orthogolan to orientation
) *Ring {
	length := int(radius * 2 * math.Pi / distPerLED)

	r := &Ring{
		Basic:  node.NewBasic(),
		Radius: radius,
		Object: space.NewObject(center, orientation, rotation),
	}
	r.Row = NewRow(length, r.getLights)
	return r
}

func (r Ring) getLights() []ifaces.Light {

	rotationMatrix := r.GetRotation().RotationMatrix()
	translationMatrix := r.GetLocation().TranslationMatrix()
	transformationMatrix := translationMatrix.Multiply(rotationMatrix)

	radPerLED := distPerLED / r.Radius

	lights := make([]ifaces.Light, r.Length)
	for i := range lights {

		phi := radPerLED * float64(i)
		sin, cos := math.Sincos(phi)

		// Location of LED if Ring was in YZ-Plane with first LED on the positive Z axis
		localLocation := space.Cartesian{
			X: 0,
			Y: r.Radius * sin,
			Z: r.Radius * cos,
		}
		worldLocation := localLocation.Transform(transformationMatrix).Cartesian()

		localOrientation := space.NewSpherical(1, math.Pi/4, phi)

		worldOrientation := localOrientation.Transform(rotationMatrix).Spherical()

		lights[i] = &light.Basic{
			Position:     i,
			NumPositions: r.Length,
			Location:     worldLocation,
			Orientation:  worldOrientation,
		}
	}
	return lights

}

// GetType returns the type
func (Ring) GetType() string {
	return "NPRing"
}
