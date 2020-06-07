package renderer

import "github.com/hajimehoshi/ebiten"

func drawImg(screen, image *ebiten.Image, x float64, y float64, scale float64) {
	pos := ebiten.GeoM{}
	pos.Scale(scale, scale)
	pos.Translate(x, y)
	screen.DrawImage(image, &ebiten.DrawImageOptions{GeoM: pos})
}
