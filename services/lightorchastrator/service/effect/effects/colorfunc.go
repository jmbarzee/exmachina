package effects

import (
	"math"
	"time"

	"github.com/jmbarzee/domain/services/lightorchastrator/service/effect/color"
)

// ColorFunc is used by effects to provide colors.
// Often placed under Effect.color
type ColorFunc func(t time.Time) color.HSLA

// NewColorFuncStatic creates a ColorFunc which only produces onlyColor.
func NewColorFuncStatic(onlyColor color.HSLA) ColorFunc {
	return func(time.Time) color.HSLA {
		return onlyColor
	}
}

// NewColorFuncStatic creates a ColorFunc which produces shifting colors starting at colorStart.
func NewColorFuncHueShift(colorStart color.HSLA, shiftFunc ShiftFunc) ColorFunc {
	return func(t time.Time) color.HSLA {
		newColor := colorStart
		newColor.ShiftHue(shiftFunc(t))
		return newColor
	}
}

// NewColorFuncStatic creates a ColorFunc which produces shifting colors,
// bouncing between colorStart and colorEnd,
// starting at colorStart and shifting in the direction specified by upSpectrum
func NewColorFuncHueBounce(colorStart, colorEnd color.HSLA, direction int, shiftFunc ShiftFunc) ColorFunc {
	// TODO build function and consider adding shift functionality to shiftFunc
	//var hueDistance float32
	if direction == 1 {

	} else if direction == -1 {

	} else {
		panic("NewColorFuncHueBounce expected direction to be 1 or -1")
	}
	return func(t time.Time) color.HSLA {
		newColor := colorStart
		newColor.ShiftHue(shiftFunc(t))
		return newColor
	}
}

type ShiftFunc func(t time.Time) float32

func NewShiftFuncLinear(start time.Time, shiftPerMillisecond float32) ShiftFunc {
	return func(t time.Time) float32 {
		milliseconds := float32(t.Sub(start) / time.Millisecond)
		return shiftPerMillisecond * milliseconds
	}
}

func NewShiftFuncSinusoidal(start time.Time, millisecondsPerCycle int, baseShiftPerMillisecond, rangeShift float32) ShiftFunc {
	return func(t time.Time) float32 {
		milliseconds := float32(t.Sub(start) / time.Millisecond)
		cycles := milliseconds / float32(millisecondsPerCycle)
		sin := float32(math.Sin(2 * math.Pi * float64(cycles)))
		adjustedSin := sin / 2 // adjusted sin has the correct frequency but y ranges from 0 to 1
		baseShift := float32(baseShiftPerMillisecond) * milliseconds
		return rangeShift*adjustedSin + baseShift
	}
}
