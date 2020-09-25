package effect

import (
	"time"

	"github.com/jmbarzee/dominion/services/lightorchastrator/service/color"
	"github.com/jmbarzee/dominion/services/lightorchastrator/service/shared"
)

type BasicEffect struct {
	shared.TimeSpan
	Rank int
}

func (e BasicEffect) Priotity() int { return e.Rank }

func (e BasicEffect) Render(time.Time, []shared.Light) []color.HSLA { return nil }
