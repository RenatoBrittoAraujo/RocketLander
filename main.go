package main

import (
	"os"

	"github.com/renatobrittoaraujo/rl/input"
	"github.com/renatobrittoaraujo/rl/manager"
)

func main() {
	args := os.Args[1:]
	trainMode := false
	inputMode := input.AIInput
	for _, arg := range args {
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
	if trainMode && inputMode != input.AIInput {
		panic("Training must be done with AI")
	}
	manager.RunSimulation(!trainMode, inputMode)
}
