package effect

import (
	"testing"
	"time"

	"github.com/jmbarzee/services/lightorchestrator/service/color"
	"github.com/jmbarzee/services/lightorchestrator/service/light"
	"github.com/jmbarzee/services/lightorchestrator/service/space"
	"github.com/jmbarzee/services/lightorchestrator/service/vibe/ifaces"
	helper "github.com/jmbarzee/services/lightorchestrator/service/vibe/testhelper"
)

type (
	EffectTest struct {
		Name         string
		Effect       ifaces.Effect
		IntialLights []light.Light
		Instants     []Instant
	}

	Instant struct {
		Time           time.Time
		ExpectedLights []light.Light
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

func GetLights(length int, c color.HSLA) []light.Light {
	lights := make([]light.Light, length)
	for i := range lights {
		lights[i] = &TestLight{
			Color: c,
		}
	}
	return lights
}

type TestLight struct {
	Color color.HSLA
}

// GetColor returns the color of the light
func (l TestLight) GetColor() color.HSLA {
	return l.Color
}

// SetColor changes the color of the light
func (l *TestLight) SetColor(newColor color.HSLA) {
	l.Color = newColor
}

// GetPosition returns the position of the Light (in a string)
func (l TestLight) GetPosition() int {
	return 0
}

// GetLocation returns the point in space where the Light is
func (l TestLight) GetLocation() space.Vector {
	return space.Vector{}
}

// GetOrientation returns the direction the Light points
func (l TestLight) GetOrientation() space.Orientation {
	return space.Orientation{}
}
