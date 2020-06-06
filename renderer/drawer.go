package renderer

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/sim"
)

func drawSimulation(screen *ebiten.Image) {
	calcFPS()

	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})

	groundPos := ebiten.GeoM{}
	groundPos.Translate(0, screenHeight*(1-groundSlicePercentage))
	screen.DrawImage(groundImage, &ebiten.DrawImageOptions{GeoM: groundPos})

	grassPos := ebiten.GeoM{}
	grassPos.Scale(1, float64(rocketScale()))
	grassPos.Translate(0, screenHeight*(1-groundSlicePercentage))
	screen.DrawImage(grassImage, &ebiten.DrawImageOptions{GeoM: grassPos})

	screen.DrawImage(rocketImage, rocketDrawData())

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

func calcFPS() {
	frames++
	if helpers.SubtractTimeInSeconds(lastRecordedTime, time.Now()) >= 1 {
		lastFPS = frames
		frames = 0
		lastRecordedTime = time.Now()
	}
}

func composePrint(rocket *sim.Rocket) (msg string) {
	msg = fmt.Sprintf(
		" FPS: %v\n Rocket Height: %0.2f\n Rocket Thrust: %0.2f%%\n",
		lastFPS,
		rocket.Position.Y,
		rocket.ThrustPercentage())

	msg += fmt.Sprintf(
		" Rocket Fuel: %0.2f%%\n Ignitions Remaining: %v\n",
		rocket.FuelPercentage(),
		rocket.EngineStartsRemaining)

	return
}

// RocketScale returns a scale of rocket size, ranging from (0.0, 1.0]
//
// the highest when rocket lies on the lowest position possible (0)
// and going down as rocket position is higher
func rocketScale() float32 {
	return (sim.RocketLenght * rocketScaleAdjust) / (rocket.Position.Y + (sim.RocketLenght * rocketScaleAdjust))
}

func rocketDrawData() *ebiten.DrawImageOptions {
	pos := ebiten.GeoM{}

	// Adjusts length of rocket to fit screen
	rocketAdjustedLength := sim.RocketLenght * rocketSizeAdjust

	// Gets rocket scale to represent a high flying rocket
	rs := float64(rocketScale())
	rs *= rocketSizeAdjust
	pos.Scale(rs, rs)

	// Change center of image to the middle of the rocket instead of the top left
	pos.Translate(-rocketAdjustedLength/2, -rocketAdjustedLength*5)
	// Then rotates it to the current rocket rotation
	pos.Rotate(float64(rocket.Direction - math.Pi/2))

	// Now sets the rocket position somewhere in screen

	// X position is affected by a visual 'lag' when rocket's x position changes, so the
	// watcher understands that the rocket is going left or right
	lag := rocket.Position.X - lastX
	if math.Abs(float64(lag)) > maxXLag {
		lag *= float32(maxXLag / math.Abs(float64(lag)))
	}
	xpos := screenWidth/2 + lag
	lastX += (rocket.Position.X - lastX) / 2

	// Y position increases ever more slowly as rocket increase size, it gives the impression
	// that the rocket is moving very far away from the ground, without ever leaving the screen
	ypos := screenHeight/2 + minGroundDist - rocket.Position.Y*float32(rs)
	pos.Translate(float64(xpos), float64(ypos))

	return &ebiten.DrawImageOptions{
		GeoM: pos,
	}
}
