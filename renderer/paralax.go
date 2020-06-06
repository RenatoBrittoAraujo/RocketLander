package renderer

import (
	_ "image/png" // Used to import png file by ebitenutil

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	// Holds paralaxes images
	paralax1 *ebiten.Image
	paralax2 *ebiten.Image
)

func init() {
	// Loads paralax images
	paralax1, _, _ = ebitenutil.NewImageFromFile("assets/Paralax1.png", ebiten.FilterDefault)
	paralax2, _, _ = ebitenutil.NewImageFromFile("assets/Paralax2.png", ebiten.FilterDefault)
}

func drawParallax(screen *ebiten.Image) {
	drawParallaxImg(screen, paralax2, 0.5)
	drawParallaxImg(screen, paralax1, 1.0)
}

func drawParallaxImg(screen, image *ebiten.Image, ratio float32) {
	x := rocket.Position.X

	imgWidth, _ := image.Size()
	posX := (int(x*ratio) % imgWidth) + imgWidth/2

	drawImg(screen, image, posX-imgWidth, 0)
	drawImg(screen, image, posX+imgWidth, 0)
	drawImg(screen, image, posX, 0)
}
