package renderer

import (
	"fmt"
	"image/color"
	"math"

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
	featureImage, _    = ebiten.NewImage(int(width*0.6), 80, ebiten.FilterDefault)
)

func init() {
	// Set color of rectangles
	groundImage.Fill(color.RGBA{153, 102, 0, 255})       // Brownish
	grassImage.Fill(color.RGBA{100, 240, 100, 255})      // Greenish
	backgroundImage.Fill(color.RGBA{120, 120, 240, 255}) // Blueish
	featureImage.Fill(color.White)
}

func drawSimulation(screen *ebiten.Image) {
	calcFPS()

	drawImg(screen, backgroundImage, 0, 0, 1)

	drawParallax(screen)

	drawProgressBar(screen, rocket.FuelPercentage(), screenWidth-60, 360, blueBar)

	velocity := rocket.Velocity()
	if velocity > sim.MaxLandingVelocity*20 {
		velocity = sim.MaxLandingVelocity * 20
	}
	speedBar := redBar
	if velocity <= sim.MaxLandingVelocity {
		speedBar = greenBar
	}
	velocity /= sim.MaxLandingVelocity * 20
	drawProgressBar(screen, velocity, screenWidth-30, 360, speedBar)

	drawRocket(screen, rocket)

	drawParticles(screen, rocket)

	groundPos := screenHeight * (1 - groundSlicePercentage)

	drawImg(screen, groundImage, 0, groundPos, 1)

	grassPos := ebiten.GeoM{}
	grassPos.Scale(1, 1)
	grassPos.Translate(0, groundPos)
	screen.DrawImage(grassImage, &ebiten.DrawImageOptions{GeoM: grassPos})

	if rocket.IsAscending() {
		text.Draw(screen, "ASCENTION", mplusBigFont, screenWidth/2-190, screenHeight-35, color.RGBA{255, 255, 255, 255})
	} else if sim.DetectGroundCollision(rocket) > 0 && !rocket.IsAscending() {
		success := sim.LandingScore(rocket) >= 0
		msgColor := color.RGBA{255, 90, 90, 255}
		if success {
			msgColor = color.RGBA{70, 200, 70, 255}
		}
		drawImg(screen, featureImage, screenWidth/5-65, 65, 1)
		text.Draw(screen, "Landing Score: "+fmt.Sprintf("%0.2f", sim.LandingScore(rocket)), mplusBigFont, screenWidth/2-420, 130, msgColor)
		text.Draw(screen, "Vertical Speed: "+fmt.Sprintf("%0.2f m/s", -rocket.SpeedVector.Y), mplusBigFont, screenWidth/2-420, 200+70, color.White)
		text.Draw(screen, "Horizontal Speed: "+fmt.Sprintf("%0.2f m/s", rocket.SpeedVector.X), mplusBigFont, screenWidth/2-420, 270+70, color.White)
		text.Draw(screen, "Angle: "+fmt.Sprintf("%0.2f°", rocket.Direction*180.0/math.Pi-90), mplusBigFont, screenWidth/2-420, 340+70, color.White)
		text.Draw(screen, "PRESS SPACE TO RESET", mplusBigFont, screenWidth/2-360, screenHeight-20, color.White)
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
