package effect

import (
	"testing"
	"time"

	"github.com/jmbarzee/color"
	"github.com/jmbarzee/services/lightorchestrator/service/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
	"github.com/jmbarzee/space"
)

type (
	EffectTest struct {
		Name         string
		Effect       ifaces.Effect
		IntialLights []ifaces.Light
		Instants     []Instant
	}

	Instant struct {
		Time           time.Time
		ExpectedLights []ifaces.Light
	}
)

func RunEffectTests(t *testing.T, cases []EffectTest) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			for i, instant := range c.Instants {
				actualLights := c.Effect.Render(instant.Time, c.IntialLights)
				for j, expectedLight := range instant.ExpectedLights {
					actualLight := actualLights[j]
					if !helper.ColorsEqual(expectedLight.GetColor(), actualLight.GetColor()) {
						t.Fatalf("instant %v, light %v failed:\n\tExpected: %v,\n\tActual: %v", i, j, expectedLight.GetColor(), actualLight.GetColor())
					}
				}
			}
		})
	}
}

func GetLights(length int, c color.HSL) []ifaces.Light {
	lights := make([]ifaces.Light, length)
	for i := range lights {
		lights[i] = &TestLight{
			Color: c,
		}
	}
	return lights
}

type TestLight struct {
	Color color.Color
}

// GetColor returns the color of the light
func (l TestLight) GetColor() color.Color {
	return l.Color
}

// SetColor changes the color of the light
func (l *TestLight) SetColor(newColor color.Color) {
	l.Color = newColor
}

// GetPosition returns the position of the Light (in a string)
func (l TestLight) GetPosition() (int, int) {
	return 0, 0
}

// GetLocation returns the point in space where the Light is
func (l TestLight) GetLocation() space.Cartesian {
	return space.Cartesian{}
}

// GetOrientation returns the direction the Light points
func (l TestLight) GetOrientation() space.Spherical {
	return space.Spherical{}
}
