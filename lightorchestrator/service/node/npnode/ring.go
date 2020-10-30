package neopixel

import (
	"math"

	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
)

// Ring is a representation of a ring of neopixels.
// Ring implements effect.Device
type Ring struct {
	// Row provides the implementation of effect.Allocater
	*Row

	node.Basic

	// Radius is the distance from the Center to any LED
	Radius float64

	space.Object
}

var _ node.Node = (*Ring)(nil)

// NewRing creates a new Ring
func NewRing(
	radius float64,
	center space.Vector,
	orientation space.Orientation,
	rotation space.Orientation,
) *Ring {
	length := int(radius * 2 * math.Pi / distPerLED)

	d := &Ring{
		Radius: radius,
		Object: space.NewObject(center, orientation, rotation),
		Row:    NewRow(length, r.getLights),
	}
	return d
}

func (r Ring) getLights() []light.Light {

	orientationMatrix := r.GetOrientation().RotationMatrix()

	rotationMatrix := r.GetRotation().RotationMatrix()

	translationMatrix := r.GetLocation().TranslationMatrix()

	transformationMatrix := translationMatrix.Mult(orientationMatrix)

	radPerLED := distPerLED / r.Radius

	lights := make([]light.Light, length)
	for i := range lights {

		localPhi := radPerLED * float64(i)

		// Location of LED if Ring was in XZ-Plane with first LED on the positive X axis
		sin, cos := math.Sincos(float64(localPhi))
		location := space.Vector{
			X: radius * cos,
			Y: radius * sin,
			Z: 0,
		}
		// Transform to match rotation and tilt of ring
		rotatedLocation := location.Transform(orientationMatrix)
		// Translate to be relative to origin and not center
		lightLocation := center.Translate(rotatedLocation)

		lightOrientation := orientation.Rotate(localPhi)

		lights[i] = &light.Basic{
			Position:     i,
			NumPositions: length,
			Location:     lightLocation,
			Orientation:  lightOrientation,
		}
	}
	return lights

}

// GetType returns the type
func (Ring) GetType() string {
	return "NPRing"
}
