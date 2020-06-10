package renderer

import (
	"fmt"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/renatobrittoaraujo/rl/sim"
	"golang.org/x/image/font"
)

var (
	// Holds loading message's font
	mplusBigFont font.Face

	// Holds last frame receive time
	lastFrameTime time.Time
)

func init() {
	// Sets font face
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)
	mplusBigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     100,
		Hinting: font.HintingFull,
	})
}

func composePrint(rocket *sim.Rocket) (msg string) {
	msg = fmt.Sprintf(
		" FPS: %v\n Rocket Height: %0.2f\n X Position: %0.2f\n Velocity: %0.2f m/s\n\n",
		lastFPS,
		rocket.Position.Y,
		rocket.Position.X,
		rocket.Velocity())

	msg += fmt.Sprintf(
		" Rocket Fuel: %0.2f%%\n Ignitions Remaining: %v\n Rocket Thrust: %0.2f%%",
		rocket.FuelPercentage(),
		rocket.EngineStartsRemaining,
		rocket.ThrustPercentage())

	return
}
