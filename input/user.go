package input

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	thrustChangePerSecond = 0.015
)

type user struct{}

func (user user) UpdateSim(rocket *sim.Rocket) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		thrust += thrustChangePerSecond
		if thrust > 1 {
			thrust = 1
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
		rocket.JetRight()
	}
}
