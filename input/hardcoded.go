package input

import (
	"github.com/renatobrittoaraujo/rl/sim"
)

type hardcoded struct{}

var thrust float32 = 0

func (hardcoded hardcoded) UpdateSim(rocket *sim.Rocket) {
	rocket.SetThrust(1)
}
