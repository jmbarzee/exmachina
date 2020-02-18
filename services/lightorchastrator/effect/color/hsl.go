package color

import (
	"math"
)

type (
	HSL struct {
		H, S, L float32
	}
)

// SetHue will change hue to h (with wrapping).
func (c *HSL) SetHue(h float32) {
	c.H = modOne(h)
}

// ShiftHue will shift hue by h (with wrapping).
func (c *HSL) ShiftHue(h float32) {
	c.SetHue(c.H + h)
}

// SetSaturation will change saturation to s (with bounding).
func (c *HSL) SetSaturation(s float32) {
	if s > Max {
		c.S = Max
	} else if s < Min {
		c.S = Min
	} else {
		c.S = s
	}
}

// SetLightness will change lightness to l (with bounding).
func (c *HSL) SetLightness(l float32) {
	if l > Max {
		c.L = Max
	} else if l < Min {
		c.L = Min
	} else {
		c.L = l
	}
}

func (c HSL) ToRGB() RGB {
	h := c.H
	s := c.S
	l := c.L

	if s == 0 {
		// it's gray
		return RGB{l, l, l}
	}

	var v1, v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return RGB{r, g, b}
}

func hueToRGB(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}

func BlendHSLIgnoreBlack(a, b HSL) HSL {
	if a.L <= Min {
		return b
	}
	if b.L <= Min {
		return a
	}
	var hsl HSL
	hsl.H = BlendHue(a.H, b.H)
	hsl.S = Average(a.S, b.S)
	hsl.L = Average(a.L, b.L)
	return hsl
}

func BlendHSL(a, b HSL) HSL {
	var hsl HSL
	hsl.H = BlendHue(a.H, b.H)
	hsl.S = Average(a.S, b.S)
	hsl.L = Average(a.L, b.L)
	return hsl
}

func BlendHSLWeighted(a, b HSL, weight float32) HSL {
	var hsl HSL
	hsl.H = BlendHueWeighted(a.H, b.H, weight)
	hsl.S = AverageWeighted(a.S, b.S, weight)
	hsl.L = AverageWeighted(a.L, b.L, weight)
	return hsl
}

// BlendHue will accuratly blend Hues, by finding their midpoint, and accounting for wraping.
func BlendHue(h1, h2 float32) float32 {
	var max, min float32
	if h1 > h2 {
		max, min = h1, h2
	} else {
		max, min = h2, h1
	}

	// distance between hues in both directions
	cont := max - min // contiguous
	wrap := 1 - max + min
	if cont < wrap {
		return min + cont/2
	} else {
		return modOne(max + wrap/2)
	}
}

// BlendHueWeighted will accuratly blend Hues, by finding their midpoint, and accounting for wraping.
// weight should range from 0 to 1 and is used to favor a hue.
// 0 will favor the first hue and 1 will favor the hue.
func BlendHueWeighted(h1, h2, weight float32) float32 {
	if h1 > h2 {
		// distance between hues in both directions
		cont := h1 - h2     // contiguous
		wrap := 1 - h1 + h2 // wrapped distance
		if cont < wrap {
			return modOne(h2 + cont*(1-weight))
		} else {
			return modOne(h1 + wrap*weight)
		}
	} else {
		// distance between hues in both directions
		cont := h2 - h1     // contiguous
		wrap := 1 - h2 + h1 // wrapped distance
		if cont < wrap {
			return modOne(h1 + cont*weight)
		} else {
			return modOne(h2 + wrap*(1-weight))
		}
	}
}

// BlendL
func BlendL(l1, l2 float32) float32 {
	if l1 < 0.5 && l2 >= 0.5 {
		return Average(l1, l2) - 0.25
	} else if l1 >= 0.5 && l2 < 0.5 {
		return Average(l1, l2) - 0.25
	}
	return Average(l1, l2)
}

func modOne(val float32) float32 {
	return float32(math.Mod(float64(val), 1.0))
}
