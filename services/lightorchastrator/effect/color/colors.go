package color

import (
	"math/rand"
)

var (
	Black = HSL{Min, Min, Min}
	Grey  = HSL{Min, Min, Half}
	White = HSL{Min, Min, Max}

	Red         = HSL{Min, Max, Half}
	WarmRed     = HSL{1.0 / 24, Max, Half}
	Orange      = HSL{2.0 / 24, Max, Half}
	WarmYellow  = HSL{3.0 / 24, Max, Half}
	Yellow      = HSL{4.0 / 24, Max, Half}
	CoolYellow  = HSL{5.0 / 24, Max, Half}
	YellowGreen = HSL{6.0 / 24, Max, Half}
	WarmGreen   = HSL{7.0 / 24, Max, Half}
	Green       = HSL{8.0 / 24, Max, Half}
	CoolGreen   = HSL{9.0 / 24, Max, Half}
	GreenCyan   = HSL{10.0 / 24, Max, Half}
	WarmCyan    = HSL{11.0 / 24, Max, Half}
	Cyan        = HSL{12.0 / 24, Max, Half}
	CoolCyan    = HSL{13.0 / 24, Max, Half}
	BlueCyan    = HSL{14.0 / 24, Max, Half}
	CoolBlue    = HSL{15.0 / 24, Max, Half}
	Blue        = HSL{16.0 / 24, Max, Half}
	WarmBlue    = HSL{17.0 / 24, Max, Half}
	Violet      = HSL{18.0 / 24, Max, Half}
	CoolMagenta = HSL{19.0 / 24, Max, Half}
	Magenta     = HSL{20.0 / 24, Max, Half}
	WarmMagenta = HSL{21.0 / 24, Max, Half}
	RedMagenta  = HSL{22.0 / 24, Max, Half}
	CoolRed     = HSL{23.0 / 24, Max, Half}
)

var (
	Colors = map[string]HSL{
		"Black": Black,
		"White": White,

		"Red":  Red,
		"Cyan": Cyan,
	}

	Schemes = map[string][]HSL{}
)

var (
	Min  = float32(0.0000001)
	Max  = float32(0.9999999)
	Half = float32(.5)
)

func RandRGB(colors []RGB) RGB {
	return colors[rand.Intn(len(colors)-1)+1]
}
