package effect

import (
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/shared"
)

type BasicEffect struct {
	shared.TimeSpan
	Rank int
}

func (e BasicEffect) Priotity() int { return e.Rank }

func (e BasicEffect) Render(time.Time, []light.Light) []color.HSLA { return nil }
