package renderer

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/sim"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	screenWidth           = 1000
	screenHeight          = 700
	groundSlicePercentage = 0.15            // How much the dirt area occupy of the total space
	grassSlicePercentage  = 0.03            // How much grass occupy of the dirt area
	minGroundDist         = 160             // This variable is an adjustment so the rocket touches the ground
	rocketScaleAdjust     = 10              // The higher this number, the less the rocket reduces in scale as it goes high
	rocketSizeAdjust      = 0.3             // Scales rocket size to fit screen
	maxXLag               = screenWidth / 3 // Drawn rocket position lags a little from actual position, as a visual feature
)

var (
	rocketChannel    chan *sim.Rocket
	rocket           *sim.Rocket
	lastRecordedTime time.Time
	frames           int
	lastFPS          int

	// The following variables are just dumb rectangles to represent all objects in screen
	backgroundImage, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	groundImage, _     = ebiten.NewImage(screenWidth, screenHeight*groundSlicePercentage, ebiten.FilterDefault)
	grassImage, _      = ebiten.NewImage(screenWidth, screenHeight*grassSlicePercentage, ebiten.FilterDefault)
	rocketImage, _     = ebiten.NewImage(sim.RocketLenght, sim.RocketLenght*10, ebiten.FilterDefault)

	// Holds last X position to create a lagging sensation on X axis of change
	lastX float32 = 0

	// Holds loading message's font
	mplusBigFont font.Face
)

// Game holds rendering state for game
type Game struct{}

func init() {
	rocketImage.Fill(color.White)                        // Bet you can't guess this color
	groundImage.Fill(color.RGBA{153, 102, 0, 255})       // Brownish
	grassImage.Fill(color.RGBA{100, 240, 100, 255})      // Greenish
	backgroundImage.Fill(color.RGBA{120, 120, 240, 255}) // Blueish
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)
	mplusBigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     100,
		Hinting: font.HintingFull,
	})
}

// Reset rendering
func Reset(rc chan *sim.Rocket) {
	rocketChannel = rc
}

// DrawSim start the drawing of the simulation
func DrawSim(rc chan *sim.Rocket, fps int64) {
	lastRecordedTime = time.Now()
	lastFPS = int(fps)
	rocketChannel = rc
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rocket Lander")
	ebiten.SetMaxTPS(int(fps))
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// Layout of screen required by ebiten interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Update what will be rendered to screen required by ebiten interface
func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func calcFPS() {
	frames++
	if helpers.SubtractTimeInSeconds(lastRecordedTime, time.Now()) >= 1 {
		lastFPS = frames
		frames = 0
		lastRecordedTime = time.Now()
	}
}

// Draw to screen, required by ebiten interface
// Drawn in order of priority in screen
func (g *Game) Draw(screen *ebiten.Image) {
	if rocketChannel != nil && len(rocketChannel) == cap(rocketChannel) {
		rocket = <-rocketChannel
		drawSimulation(screen)
	} else {
		drawLoadingScreen(screen)
	}
}

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

	ebitenutil.DebugPrint(screen, composePrint(rocket))
}

func drawLoadingScreen(screen *ebiten.Image) {
	text.Draw(screen, "   LOADING\nSIMULATION", mplusBigFont, screenWidth/2-200, screenHeight/2-50, color.White)
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
