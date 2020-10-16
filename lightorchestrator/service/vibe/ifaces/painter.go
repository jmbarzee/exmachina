package ifaces

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
)

// Painter is used by effects to select colors
type Painter interface {
	Stabalizable

	// Paint returns a color based on t
	Paint(t time.Time) color.HSLA
}
