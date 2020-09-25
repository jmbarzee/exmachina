package color

import (
	"math/rand"
)

var (
	Min  = float32(0.0000001)
	Max  = float32(0.9999999)
	Half = float32(.5)
)

var (
	Black        = NewHSLA(Min, Min, Min, Min)
	Grey         = NewHSLA(Min, Min, Half, Half)
	GreyNatural  = NewHSLA(Min, Min, Min, Half)
	White        = NewHSLA(Min, Min, Max, Max)
	WhiteNatural = NewHSLA(Min, Min, Min, Max)

	Red         = NewHSLA(Min, Max, Half, Min)
	WarmRed     = NewHSLA(1.0/24, Max, Half, Min)
	Orange      = NewHSLA(2.0/24, Max, Half, Min)
	WarmYellow  = NewHSLA(3.0/24, Max, Half, Min)
	Yellow      = NewHSLA(4.0/24, Max, Half, Min)
	CoolYellow  = NewHSLA(5.0/24, Max, Half, Min)
	YellowGreen = NewHSLA(6.0/24, Max, Half, Min)
	WarmGreen   = NewHSLA(7.0/24, Max, Half, Min)
	Green       = NewHSLA(8.0/24, Max, Half, Min)
	CoolGreen   = NewHSLA(9.0/24, Max, Half, Min)
	GreenCyan   = NewHSLA(10.0/24, Max, Half, Min)
	WarmCyan    = NewHSLA(11.0/24, Max, Half, Min)
	Cyan        = NewHSLA(12.0/24, Max, Half, Min)
	CoolCyan    = NewHSLA(13.0/24, Max, Half, Min)
	BlueCyan    = NewHSLA(14.0/24, Max, Half, Min)
	CoolBlue    = NewHSLA(15.0/24, Max, Half, Min)
	Blue        = NewHSLA(16.0/24, Max, Half, Min)
	WarmBlue    = NewHSLA(17.0/24, Max, Half, Min)
	Violet      = NewHSLA(18.0/24, Max, Half, Min)
	CoolMagenta = NewHSLA(19.0/24, Max, Half, Min)
	Magenta     = NewHSLA(20.0/24, Max, Half, Min)
	WarmMagenta = NewHSLA(21.0/24, Max, Half, Min)
	RedMagenta  = NewHSLA(22.0/24, Max, Half, Min)
	CoolRed     = NewHSLA(23.0/24, Max, Half, Min)
)

var (
	Colors = map[string]HSLA{
		"Black": Black,
		"White": White,

		"Red":  Red,
		"Cyan": Cyan,
	}

	Schemes = map[string][]HSLA{}
)

func RandRGB(colors []RGB) RGB {
	return colors[rand.Intn(len(colors)-1)+1]
}
