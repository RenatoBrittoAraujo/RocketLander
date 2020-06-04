package input

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	thrustChangePerSecond = 1.5
)

type user struct{}

var thrust float32 = 0

func (user user) UpdateSim(rocket *sim.Rocket) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		thrust += thrustChangePerSecond
		if thrust > 100 {
			thrust = 100
		}
		rocket.SetThrust(thrust)
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		thrust -= thrustChangePerSecond
		if thrust < 0 {
			thrust = 0
		}
		rocket.SetThrust(thrust)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rocket.JetLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rocket.JetLeft()
	}
}
