package sim

import "math"

const (
	// RocketLenght is 70 meters (Falcon 9)
	RocketLenght = 70
)

// Point represents 2D point with X and Y
type Point struct {
	X, Y float32
}

// Vector represents 2D vector with X and Y
type Vector struct {
	X, Y float32
}

// GameState holds all relevant game data
type GameState struct {
	// RocketPosition given with {X, Y}
	RocketPosition Point
	// RocketDirection given in radians (0 is vertical up)
	RocketDirection float32
	// Speed given by vector of meters per second
	SpeedVector Vector
	Fuel        float32
}

// CreateGameState creates and returns an instace of gamestate
func CreateGameState() (gameState *GameState) {
	return &GameState{}
}

var lowg float32 = 0
var back bool = false

// Update the gamestate to it's next frame following physics
func Update(gameState *GameState) {
	lowg += 3.1415 / 50
	gameState.RocketPosition.Y = 100 + float32(math.Sin(float64(lowg)))*100
	if gameState.RocketPosition.X < 100 {
		back = false
	} else if gameState.RocketPosition.X > 10000 {
		back = true
	}
	if back {
		gameState.RocketPosition.X -= (gameState.RocketPosition.X + 10) / 100
	} else {
		gameState.RocketPosition.X += (gameState.RocketPosition.X + 10) / 50
	}
}
