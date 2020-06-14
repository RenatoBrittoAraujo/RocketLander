package main

import (
	"os"
	"strconv"

	"github.com/renatobrittoaraujo/rl/appmanager"
	"github.com/renatobrittoaraujo/rl/input"
)

func main() {
	args := os.Args[1:]
	trainMode := false
	inputMode := input.AIInput
	var seed, fps int
	for _, arg := range args {
		if arg[0:3] == "fps" {
			fps, _ = strconv.Atoi(arg[4:])
			continue
		}
		if arg[0:4] == "seed" {
			seed, _ = strconv.Atoi(arg[5:])
			continue
		}
		switch arg {
		case "train":
			trainMode = true
		case "hardcoded":
			inputMode = input.HardcodedInput
		case "user":
			inputMode = input.UserInput
		case "ai":
			inputMode = input.AIInput
		case "draw":
			trainMode = false
		default:
			panic("Invalid CLI argument")
		}
	}
	appmanager.StartSimulationDriver(!trainMode, inputMode, seed, fps)
}
