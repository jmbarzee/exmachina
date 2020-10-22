package color

import "math"

type RGBA struct {
	RGB
	A float64
}

func (c RGBA) ToUInt32WGRB() uint32 {
	val := uint32(c.B*math.MaxUint8) << 0
	val |= uint32(c.R*math.MaxUint8) << 8
	val |= uint32(c.G*math.MaxUint8) << 16
	val |= uint32(c.A*math.MaxUint8) << 24
	return val
}
func FromUInt32WGRB(wgrb uint32) RGBA {
	mask := uint32(0x000000ff)

	uint8b := mask & (wgrb >> 0)
	b := float64(uint8b) / math.MaxFloat64
	uint8r := mask & (wgrb >> 8)
	r := float64(uint8r) / math.MaxFloat64
	uint8g := mask & (wgrb >> 16)
	g := float64(uint8g) / math.MaxFloat64
	uint8a := mask & (wgrb >> 24)
	a := float64(uint8a) / math.MaxFloat64
	return RGBA{
		RGB: RGB{
			R: r,
			G: g,
			B: b,
		},
		A: a,
	}
}
