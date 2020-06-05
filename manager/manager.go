package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
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
var rocketChannel chan *sim.Rocket

// StartSimulationDriver receives draw bool and run sim with screen drawing or on CLI
func StartSimulationDriver(draw bool, inputType int) {
	rocket := sim.CreateRocket()
	inputManager, err := input.CreateInput(inputType)
	rocketChannel = nil
	if err {
		panic("Input \"" + input.InputString[inputType] + "\" has not been initalized correctly")
	}
	wg.Add(1)
	go func() {
		for range time.Tick(time.Second / time.Duration(60)) {
			if ebiten.IsKeyPressed(ebiten.KeyR) {
				wg.Done()
				break
			}
		}
	}()
	if draw {
		rocketChannel = make(chan *sim.Rocket)
		go startSimulation(rocket, simDrawFrames*5, inputManager)
		renderer.DrawSim(rocketChannel, simDrawFrames)
		wg.Done()
	} else {
		startSimulation(rocket, simCliFrames, inputManager)
		wg.Done()
	}
	wg.Wait()
	close(rocketChannel)
	StartSimulationDriver(draw, inputType)
}

func startSimulation(rocket *sim.Rocket, fps int64, inputManager input.Manager) {
	var frames int64 = 0
	for range time.Tick(time.Second / time.Duration(fps)) {
		if frames%(fps*5) == 0 {
			fmt.Println("Simulation Frames Calculated:", frames)
		}
		frames++
		if rocket.IsAscending() && false {
			rocket.Ascend(1)
		} else {
			inputManager.UpdateSim(rocket)
		}
		rocket.Update()
		if col := sim.DetectGroundCollision(rocket); col > 0 && !rocket.IsAscending() {
			close(rocketChannel)
			wg.Done()
			break
		}
		if rocketChannel != nil {
			go emmitRocketState(rocket)
		}
	}
}

func emmitRocketState(rocket *sim.Rocket) {
	rocketChannel <- rocket
}
