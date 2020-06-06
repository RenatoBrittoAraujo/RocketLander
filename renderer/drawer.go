package renderer

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	groundSlicePercentage = 0.15            // How much the dirt area occupy of the total space
	grassSlicePercentage  = 0.03            // How much grass occupy of the dirt area
	minGroundDist         = 168             // This variable is an adjustment so the rocket touches the ground
	rocketScaleAdjust     = 10              // The higher this number, the less the rocket reduces in scale as it goes high
	rocketSizeAdjust      = 0.3             // Scales rocket size to fit screen
	maxXLag               = screenWidth / 3 // Drawn rocket position lags a little from actual position, as a visual feature
)

var (
	height float32 = screenHeight
	width  float32 = screenWidth

	// The following variables are just dumb rectangles to represent all objects in screen
	backgroundImage, _ = ebiten.NewImage(int(width), int(height), ebiten.FilterDefault)
	groundImage, _     = ebiten.NewImage(int(width), int(height*groundSlicePercentage), ebiten.FilterDefault)
	grassImage, _      = ebiten.NewImage(int(width), int(height*grassSlicePercentage), ebiten.FilterDefault)
	rocketImage, _     = ebiten.NewImage(sim.RocketLenght, sim.RocketLenght*10, ebiten.FilterDefault)

	// Holds last X position to create a lagging sensation on X axis of change
	lastX float32 = 0
)

func init() {
	// Set color of rectangles
	rocketImage.Fill(color.White)                        // Bet you can't guess this color
	groundImage.Fill(color.RGBA{153, 102, 0, 255})       // Brownish
	grassImage.Fill(color.RGBA{100, 240, 100, 255})      // Greenish
	backgroundImage.Fill(color.RGBA{120, 120, 240, 255}) // Blueish

}

func drawSimulation(screen *ebiten.Image) {
	calcFPS()

	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})

	drawParallax(screen)

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
