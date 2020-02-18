package effects

import "github.com/jmbarzee/domain/services/lightorchastrator/effect/shared"

type BasicEffect struct {
	shared.TimeSpan
	Rank int
}

func (e BasicEffect) Priotity() int { return e.Rank }

//func (e BasicEffect) Render(time.Time, []effect.Light) []color.HSLA { return nil }
