package manager

import (
	"fmt"
	"sync"
	"time"

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
	gameState := sim.CreateGameState()
	wg.Add(1)
	if draw {
		go startSimUpdater(gameState, simDrawFrames)
		startDrawer(gameState, simDrawFrames)
	} else {
		go startSimUpdater(gameState, simCliFrames)
	}
	wg.Wait()
}

func startDrawer(gameState *sim.GameState, fps int64) {
	defer wg.Done()
	renderer.DrawSim(gameState, fps)
}

func startSimUpdater(gameState *sim.GameState, fps int64) {
	var frames int64 = 0
	for range time.Tick(time.Second / time.Duration(fps)) {
		if frames%(fps*3) == 0 {
			fmt.Println("Simulation Frames Calculated:", frames)
		}
		frames++
		sim.Update(gameState)
	}
}
