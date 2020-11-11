package npdevice

import (
	"github.com/jmbarzee/services/lightorchestrator/service/device"
	"github.com/jmbarzee/services/lightorchestrator/service/node"
	"github.com/jmbarzee/services/lightorchestrator/service/node/npnode"
	"github.com/jmbarzee/space"
)

const (
	npBarLength  = 2
	ledsPerMeter = 60

	ledsPerNPBar = npBarLength * ledsPerMeter
)

// Bar is a strait bar of lights
type Bar struct {
	device.Basic

	*npnode.Line
}

var _ device.Device = (*Bar)(nil)

// NewBar creates a new Bar
func NewBar(uuid string, start space.Cartesian, direction, rotation space.Spherical) Bar {
	return Bar{
		Basic: device.NewBasic(uuid),
		Line:  npnode.NewLine(ledsPerNPBar, start, direction, rotation),
	}
}

// GetNodes returns all the Nodes which the device holds
func (b Bar) GetNodes() []node.Node {
	return []node.Node{
		b.Line,
	}
}

// GetType returns the type
func (Bar) GetType() string {
	return "npBar"
}
