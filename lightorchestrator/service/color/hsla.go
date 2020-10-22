package color

type (
	HSLA struct {
		HSL
		A float64
	}
)

func NewHSLA(h, s, l, a float64) HSLA {
	return HSLA{
		HSL: HSL{
			H: h,
			S: s,
			L: l,
		},
		A: a,
	}
}

func (c HSLA) ToRGBA() RGBA {
	return RGBA{
		RGB: c.ToRGB(),
		A:   c.A,
	}
}
