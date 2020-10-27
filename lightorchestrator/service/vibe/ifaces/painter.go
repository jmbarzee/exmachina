package ifaces

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
)

// Painter is used by effects to select colors
type Painter interface {
	Stabalizable

	// Paint returns a color based on t
	Paint(t time.Time, l light.Light) color.HSLA
}
