package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/renatobrittoaraujo/rl/input"
	"github.com/renatobrittoaraujo/rl/renderer"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	// simDrawFrames is the amount of frames per second a drawn simulation has
	simDrawFrames = 60
	// simCliFrames is the max amount of interations of simulation per second in CLI simulation
	simCliFrames = 1000
)

var wg sync.WaitGroup

// RunSimulation receives draw bool and run sim with screen drawing or on CLI
func RunSimulation(draw bool, inputType int) {
	rocket := sim.CreateRocket()
	inputManager, err := input.CreateInput(inputType)
	rocketChannel := make(chan *sim.Rocket)
	if err {
		panic("Input \"" + input.InputString[inputType] + "\" has not been initalized correctly")
	}
	wg.Add(1)
	if draw {
		go startSimUpdater(rocket, simDrawFrames, inputManager, rocketChannel)
		renderer.DrawSim(rocketChannel, simDrawFrames)
		wg.Done()
	} else {
		startSimUpdater(rocket, simCliFrames, inputManager, rocketChannel)
		wg.Done()
	}
	wg.Wait()
}

func startSimUpdater(rocket *sim.Rocket, fps int64, inputManager input.Manager, channel chan *sim.Rocket) {
	var frames int64 = 0
	for range time.Tick(time.Second / time.Duration(fps)) {
		if frames%(fps*5) == 0 {
			fmt.Println("Simulation Frames Calculated:", frames)
		}
		frames++
		inputManager.UpdateSim(rocket)
		sim.UpdateRocket(rocket)
		go emmitRocketState(channel, rocket)
	}
}

func emmitRocketState(channel chan *sim.Rocket, rocket *sim.Rocket) {
	channel <- rocket
}
