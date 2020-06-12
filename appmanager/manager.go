package appmanager

import (
	"fmt"
	"math/rand"
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

var (
	rocketChannel chan *sim.Rocket
	inputType     int
	draw          bool
	loaded        bool
	inputManager  input.Manager
)

// StartSimulationDriver receives draw bool and run sim with screen drawing or on CLI
func StartSimulationDriver(argdraw bool, arginputType int) {
	inputType = arginputType
	draw = argdraw
	linputManager, err := input.CreateInput(inputType)
	if err {
		panic("Input \"" + input.InputString[inputType] + "\" has not been initalized correctly")
	}
	inputManager = linputManager
	rocketChannel = make(chan *sim.Rocket, 1)
	go startSimulationInstance()
	renderer.DrawSim(rocketChannel, simDrawFrames)
}

func startSimulationInstance() {
	for {
		rocket := sim.CreateRocket()
		rand.Seed(time.Now().Local().UnixNano())
		seed := rand.Int()*100000000 - 50000000
		fmt.Println("SEED USED:", seed)
		var fps int
		if draw {
			fps = simDrawFrames * 20
		} else {
			fps = simCliFrames
		}
		if draw {
			waitKeyPress(ebiten.KeySpace, nil)
		}
		for range time.Tick(time.Second / time.Duration(fps)) {
			if draw && ebiten.IsKeyPressed(ebiten.KeyR) {
				break
			}
			if rocket.IsAscending() {
				rocket.Ascend(float32(seed))
			} else {
				inputManager.UpdateSim(rocket)
			}
			rocket.Update()
			if col := sim.DetectGroundCollision(rocket); col > 0 && !rocket.IsAscending() {
				if draw {
					waitKeyPress(ebiten.KeySpace, rocket)
					time.Sleep(time.Millisecond * 70 /* When spacebar is pressed at the end of simulation, little lag so no overlap with next spacebar press */)
				}
				break
			}
			if len(rocketChannel) < cap(rocketChannel) {
				rocketChannel <- rocket
			}
		}
	}
}

func waitKeyPress(key ebiten.Key, rocket *sim.Rocket) {
	if rocket != nil {
		for range time.Tick(time.Second / simDrawFrames) {
			if len(rocketChannel) < cap(rocketChannel) {
				rocketChannel <- rocket
			}
			if ebiten.IsKeyPressed(key) {
				break
			}
		}
	} else {
		for range time.Tick(time.Second / simDrawFrames) {
			if ebiten.IsKeyPressed(key) {
				break
			}
		}
	}
}
