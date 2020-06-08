package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	groundSlicePercentage = 0.15 // How much the dirt area occupy of the total space
	grassSlicePercentage  = 0.03 // How much grass occupy of the dirt area
)

var (
	height float32 = screenHeight
	width  float32 = screenWidth

	// The following variables are just dumb rectangles to represent all objects in screen
	backgroundImage, _ = ebiten.NewImage(int(width), int(height), ebiten.FilterDefault)
	groundImage, _     = ebiten.NewImage(int(width), int(height*groundSlicePercentage), ebiten.FilterDefault)
	grassImage, _      = ebiten.NewImage(int(width), int(height*grassSlicePercentage), ebiten.FilterDefault)
)

func init() {
	// Set color of rectangles
	groundImage.Fill(color.RGBA{153, 102, 0, 255})       // Brownish
	grassImage.Fill(color.RGBA{100, 240, 100, 255})      // Greenish
	backgroundImage.Fill(color.RGBA{120, 120, 240, 255}) // Blueish

}

func drawSimulation(screen *ebiten.Image) {
	calcFPS()

	drawImg(screen, backgroundImage, 0, 0, 1)

	drawParallax(screen)

	groundPos := screenHeight * (1 - groundSlicePercentage)

	drawImg(screen, groundImage, 0, groundPos, 1)

	grassPos := ebiten.GeoM{}
	grassPos.Scale(1, 1)
	grassPos.Translate(0, groundPos)
	screen.DrawImage(grassImage, &ebiten.DrawImageOptions{GeoM: grassPos})

	drawRocket(screen, rocket)
	drawParticles(screen, rocket)

	if sim.DetectGroundCollision(rocket) > 0 && !rocket.IsAscending() {
		text.Draw(screen, "COLLISION DETECTED!", mplusBigFont, screenWidth/2-360, screenHeight/2-150, color.RGBA{255, 70, 70, 255})
		text.Draw(screen, "PRESS SPACE TO RESET", mplusBigFont, screenWidth/2-360, screenHeight/2-50, color.White)
	} else {
		ebitenutil.DebugPrint(screen, composePrint(rocket))
	}
}

func drawLoadingScreen(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})

	groundPos := ebiten.GeoM{}
	groundPos.Translate(0, screenHeight*(1-groundSlicePercentage))
	screen.DrawImage(groundImage, &ebiten.DrawImageOptions{GeoM: groundPos})

	grassPos := ebiten.GeoM{}
	grassPos.Translate(0, screenHeight*(1-groundSlicePercentage))
	screen.DrawImage(grassImage, &ebiten.DrawImageOptions{GeoM: grassPos})

	text.Draw(screen, "PRESS SPACE\n  TO START", mplusBigFont, screenWidth/2-200, screenHeight/2-50, color.White)
}
