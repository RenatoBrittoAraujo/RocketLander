package input

import "github.com/renatobrittoaraujo/rl/sim"

const (
	// UserInput is a signal to Input Package that the input for current program is from the user
	UserInput = iota
	// AIInput is a signal to Input Package that the input for current program is from an AI algorithm
	AIInput
	// HardcodedInput is a signal to Input Package that the input for current program is from a hardcoded algorithm
	HardcodedInput
)

// InputString is the name of input type for a given input value
var InputString [3]string = [3]string{"User", "AI", "Hardcoded"}

// Manager is a interface that allows for easy interaction with input type
type Manager interface {
	UpdateSim(*sim.Rocket)
}

// CreateInput returns a struct that follows input.Manager interface
func CreateInput(inputType int) (Manager, bool) {
	switch inputType {
	case UserInput:
		return user{}, false
	// case AIInput:
	// 	return ai{}
	// case HardcodedInput:
	// 	return hardcoded{}
	default:
		return nil, true
	}
}

// type ai struct{}

// type hardcoded struct{}
