package renderer

import (
	"log"
	"time"

	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/sim"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth        = 1280
	screenHeight       = 720
	frameFeedTolerance = 0.500 // Seconds for how much time module waits after stopping receiving frames to declare end of simulation
)

var (
	rocketChannel    chan *sim.Rocket
	rocket           *sim.Rocket
	lastRecordedTime time.Time
	frames           int
	lastFPS          int
)

// Game holds rendering state for game
type Game struct{}

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

func calcFPS() {
	frames++
	if helpers.SubtractTimeInSeconds(lastRecordedTime, time.Now()) >= 1 {
		lastFPS = frames
		frames = 0
		lastRecordedTime = time.Now()
	}
}
