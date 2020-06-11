package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

const (
	barWidth  = 20
	barHeight = 250
)

var (
	emptyBar, _ = ebiten.NewImage(barWidth, barHeight, ebiten.FilterDefault)
	redBar, _   = ebiten.NewImage(barWidth, barHeight, ebiten.FilterDefault)
	greenBar, _ = ebiten.NewImage(barWidth, barHeight, ebiten.FilterDefault)
	blueBar, _  = ebiten.NewImage(barWidth, barHeight, ebiten.FilterDefault)
)

func init() {
	emptyBar.Fill(color.RGBA{110, 110, 110, 150})
	redBar.Fill(color.RGBA{255, 130, 130, 255})
	blueBar.Fill(color.RGBA{130, 160, 210, 255})
	greenBar.Fill(color.RGBA{170, 250, 170, 255})
}

func drawProgressBar(screen *ebiten.Image, percent float32, xpos int, ypos int, bar *ebiten.Image) {
	posempty := ebiten.GeoM{}
	posempty.Translate(float64(xpos), float64(ypos))
	screen.DrawImage(emptyBar, &ebiten.DrawImageOptions{GeoM: posempty})
	posfilled := ebiten.GeoM{}
	posfilled.Scale(1, float64(percent))
	posfilled.Translate(float64(xpos), float64(float32(ypos)+(1.0-percent)*barHeight))
	screen.DrawImage(bar, &ebiten.DrawImageOptions{GeoM: posfilled})
}
