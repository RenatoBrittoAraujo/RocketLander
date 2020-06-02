package input

// "github.com/renatobrittoaraujo/rl/sim"

const (
	// UserInput is a signal to Input Package that the input for current program is from the user
	UserInput = iota
	// AIInput is a signal to Input Package that the input for current program is from an AI algorithm
	AIInput
	// HardcodedInput is a signal to Input Package that the input for current program is from a hardcoded algorithm
	HardcodedInput
)

// // Manager is a interface that allows for easy interaction with input type
// type Manager interface {
// 	updateSim(gameState *sim.GameState) err bool
// }

// // CreateInput returns a struct that follows Manager interface
// func CreateInput(inputType int) (input Manager, err bool) {
// 	switch inputType {
// 	case UserInput:

// 	case AIInput:

// 	case HardcodedInput:

// 	default:
// 		return nil, true
// 	}
// }

// type user struct{}

// type ai struct{}

// type hardcoded struct{}

// func (user user) updateSim(gameState *GameState) err bool {
// 	err = false
// 	if ebiten.IsKeyPressed(ebiten.KeyUp) {
// 		fmt.Println("Up key pressed")
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyDown) {
// 		fmt.Println("Down key pressed")
// 	}
// 	return gameState
// }
