package npnode

import (
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/space"
)

// Line is a representation of a strait line of neopixels.
type Line struct {
	// Row provides the implementation of effect.Allocater
	*Row

	node.Basic

	*space.Object
}

var _ node.Node = (*Line)(nil)

// NewLine creates a new Line
func NewLine(
	length int,
	start space.Cartesian,
	orientation space.Spherical,
	rotation space.Spherical,
) *Line {

	l := &Line{
		Basic:  node.NewBasic(),
		Object: space.NewObject(start, orientation, rotation),
	}
	l.Row = NewRow(length, l.getLights)

	return l
}

func (l *Line) getLights() []ifaces.Light {

	singleLEDVector := l.GetOrientation()
	singleLEDVector.R = distPerLED

	lights := make([]ifaces.Light, l.Length)
	for i := range lights {
		thisLEDVector := singleLEDVector.Scale(float64(i))
		lightLocation := l.GetLocation().Translate(thisLEDVector).Cartesian()
		lightOrientation := l.GetRotation()
		lights[i] = &light.Basic{
			Position:     i,
			NumPositions: l.Length,
			Location:     lightLocation,
			Orientation:  lightOrientation,
		}
	}
	return lights
}

// GetType returns the type
func (Line) GetType() string {
	return "NPLine"
}
