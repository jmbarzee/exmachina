package shared

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
)

// Painter is used by effects to provide colors.
// Often placed under Effect.color
type Painter interface {
	// GetColor returns a color based on t
	GetColor(t time.Time) color.HSLA
	stabilizable
}
