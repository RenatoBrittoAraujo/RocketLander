package renderer

import (
	_ "image/png" // Used to import png file by ebitenutil

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	paralaxMinSize      = 0.5 // Percent
	changeAsHeightGrows = 700
)

var (
	// Holds paralaxes images
	paralax1 *ebiten.Image
	paralax2 *ebiten.Image
	paralax3 *ebiten.Image
)

func init() {
	// Loads paralax images
	paralax1, _, _ = ebitenutil.NewImageFromFile("assets/Paralax1.png", ebiten.FilterDefault)
	paralax2, _, _ = ebitenutil.NewImageFromFile("assets/Paralax2.png", ebiten.FilterDefault)
	paralax3, _, _ = ebitenutil.NewImageFromFile("assets/Paralax3.png", ebiten.FilterDefault)
}

func drawParallax(screen *ebiten.Image) {
	drawParallaxImg(screen, paralax3, 0.2)
	drawParallaxImg(screen, paralax2, 0.5)
	drawParallaxImg(screen, paralax1, 1.0)
}

func drawParallaxImg(screen, image *ebiten.Image, ratio float32) {
	x := rocket.Position.X
	y := rocket.Position.Y - sim.RocketLenght/2
	scale := paralaxScale(y, ratio)

	imgWidth, _ := image.Size()
	imgWidth = int(float64(imgWidth) * scale)
	posX := float64((int(x*ratio) % imgWidth) + imgWidth/2)
	posY := (1.0-scale)*float64(screenHeight)*1.2 - groundSlicePercentage*screenHeight
	posX += (1.0 - scale) * float64(imgWidth) / 2

	drawImg(screen, image, posX-2*float64(imgWidth), posY, scale)
	drawImg(screen, image, posX-float64(imgWidth), posY, scale)
	drawImg(screen, image, posX, posY, scale)
	drawImg(screen, image, posX+float64(imgWidth), posY, scale)
}

func paralaxScale(height, ratio float32) float64 {
	scaling := height*ratio/changeAsHeightGrows + 2 /* because height=0 -> scale=1 always */
	return float64(1/scaling + paralaxMinSize)
}
