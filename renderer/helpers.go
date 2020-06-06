package renderer

import "github.com/hajimehoshi/ebiten"

func drawImg(screen, image *ebiten.Image, x int, y int) {
	pos := ebiten.GeoM{}
	pos.Translate(float64(x), float64(y))
	screen.DrawImage(image, &ebiten.DrawImageOptions{GeoM: pos})
}
