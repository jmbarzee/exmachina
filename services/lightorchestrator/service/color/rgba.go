package color

import "math"

type RGBA struct {
	RGB
	A float32
}

func (c RGBA) ToUInt32WGRB() uint32 {
	val := uint32(c.B*math.MaxUint8) << 0
	val |= uint32(c.R*math.MaxUint8) << 8
	val |= uint32(c.G*math.MaxUint8) << 16
	val |= uint32(c.A*math.MaxUint8) << 24
	return val
}
