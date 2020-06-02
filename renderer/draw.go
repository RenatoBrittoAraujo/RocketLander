package renderer

import (
	"image/color"
	"log"

	"github.com/renatobrittoaraujo/rl/sim"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth           = 1000
	screenHeight          = 700
	groundSlicePercentage = 0.15
)

var (
	gameState          *sim.GameState
	rocketImage, _     = ebiten.NewImage(sim.RocketLenght, sim.RocketLenght*10, ebiten.FilterDefault)
	backgroundImage, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	groundImage, _     = ebiten.NewImage(screenWidth, screenHeight*groundSlicePercentage, ebiten.FilterDefault)
)

// Game holds rendering state for game
type Game struct{}

func init() {
	rocketImage.Fill(color.White)
	groundImage.Fill(color.RGBA{153, 102, 0, 255})
	backgroundImage.Fill(color.RGBA{120, 120, 240, 255})
}

// DrawSim start the drawing of the simulation
func DrawSim(gs *sim.GameState, fps int64) {
	gameState = gs
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rocket Lander")
	ebiten.SetMaxTPS(int(fps))
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

// Draw to screen, required by ebiten interface
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImage, &ebiten.DrawImageOptions{})
	groundPos := ebiten.GeoM{}
	groundPos.Translate(0, screenHeight*(1-groundSlicePercentage))
	screen.DrawImage(groundImage, &ebiten.DrawImageOptions{GeoM: groundPos})
	screen.DrawImage(rocketImage, rocketDrawData())
}

// RocketScale returns a scale of rocket size, ranging from (0.0, 1.0]
//
// the highest when rocket lies on the lowest position possible (0)
// and going down as rocket position is higher
func rocketScale() float32 {
	return (sim.Height * 3.4) / (gameState.RocketPosition.Y + (sim.Height * 3.4))
}

var rot float32
var vpos float32 = 10

func rocketDrawData() *ebiten.DrawImageOptions {
	pos := ebiten.GeoM{}
	rot = rot + 3.1415/100
	rs := float64(rocketScale())
	pos.Scale(rs, rs)
	pos.Rotate(float64(rot))
	pos.Translate(screenWidth/2, float64(screenHeight/2-gameState.RocketPosition.Y))
	return &ebiten.DrawImageOptions{
		GeoM: pos,
	}
}
