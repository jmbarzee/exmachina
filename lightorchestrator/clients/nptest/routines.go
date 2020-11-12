package nptest

import (
	"context"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/jmbarzee/color"
)

func (l *NPTest) updateLights(ctx context.Context, t time.Time) {
	// Advance the light plan
	next := l.LightPlan.Advance(t)
	if next != nil {
		imd := imdraw.New(nil)
		for i, wrgb := range next.Lights {
			rgba := color.FromUInt32WRGB(wrgb)
			x := float64(i * pixelsPerLight)
			y := 0.0

			// draw the colored lights
			imd.Color = pixel.RGB(float64(rgba.R), float64(rgba.G), float64(rgba.B))
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y))
			imd.Polygon(0)

			y = pixelsPerLight
			// draw the white lights
			imd.Color = pixel.RGB(float64(rgba.A), float64(rgba.A), float64(rgba.A))
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y))
			imd.Polygon(0)

		}
		for i, c := range color.AllColors {
			rgb := c.RGB()
			x := float64(i * pixelsPerLight)
			y := pixelsPerLight * 2.0

			// draw the colored lights
			imd.Color = pixel.RGB(rgb.R, rgb.G, rgb.B)
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y+pixelsPerLight))
			imd.Push(pixel.V(x+pixelsPerLight, y))
			imd.Polygon(0)

		}
		if !l.Window.Closed() {
			imd.Draw(l.Window)
			l.Window.Update()
		}
	}
}
