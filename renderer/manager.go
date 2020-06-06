package renderer

import (
	"image/color"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/sim"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
)

const (
	screenWidth           = 1280
	screenHeight          = 720
	groundSlicePercentage = 0.15            // How much the dirt area occupy of the total space
	grassSlicePercentage  = 0.03            // How much grass occupy of the dirt area
	minGroundDist         = 168             // This variable is an adjustment so the rocket touches the ground
	rocketScaleAdjust     = 10              // The higher this number, the less the rocket reduces in scale as it goes high
	rocketSizeAdjust      = 0.3             // Scales rocket size to fit screen
	maxXLag               = screenWidth / 3 // Drawn rocket position lags a little from actual position, as a visual feature
	frameFeedTolerance    = 0.500           // Seconds for how much time module waits after stopping receiving frames to declare end of simulation
)

var (
	rocketChannel    chan *sim.Rocket
	rocket           *sim.Rocket
	lastRecordedTime time.Time
	frames           int
	lastFPS          int

	// The following variables are just dumb rectangles to represent all objects in screen
	height             float32 = screenHeight
	width              float32 = screenWidth
	backgroundImage, _         = ebiten.NewImage(int(width), int(height), ebiten.FilterDefault)
	groundImage, _             = ebiten.NewImage(int(width), int(height*groundSlicePercentage), ebiten.FilterDefault)
	grassImage, _              = ebiten.NewImage(int(width), int(height*grassSlicePercentage), ebiten.FilterDefault)
	rocketImage, _             = ebiten.NewImage(sim.RocketLenght, sim.RocketLenght*10, ebiten.FilterDefault)

	// Holds last X position to create a lagging sensation on X axis of change
	lastX float32 = 0

	// Holds loading message's font
	mplusBigFont font.Face

	// Holds last frame receive time
	lastFrameTime time.Time
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

// DrawSim start the drawing of the simulation
func DrawSim(rc chan *sim.Rocket, fps int64) {
	lastRecordedTime = time.Now()
	lastFPS = int(fps)
	rocketChannel = rc
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rocket Lander")
	ebiten.SetMaxTPS(int(fps))
	ebiten.SetRunnableOnUnfocused(true)
	lastFrameTime = time.Now()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// Draw to screen, required by ebiten interface
// Drawn in order of priority in screen
func (g *Game) Draw(screen *ebiten.Image) {
	if rocketChannel != nil && len(rocketChannel) == cap(rocketChannel) {
		rocket = <-rocketChannel
		drawSimulation(screen)
		lastFrameTime = time.Now()
	} else if rocket != nil && helpers.SubtractTimeInSeconds(lastFrameTime, time.Now()) < frameFeedTolerance {
		drawSimulation(screen)
	} else {
		drawLoadingScreen(screen)
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
