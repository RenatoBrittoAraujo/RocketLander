package sim

import "math"

const (
	// RocketLenght is the lenght of the rocket
	RocketLenght = 20
	Height       = 100
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
	RocketPosition  Point
	RocketDirection Vector
	SpeedVector     Vector
	Fuel            float32
}

// CreateGameState creates and returns an instace of gamestate
func CreateGameState() (gameState *GameState) {
	return &GameState{}
}

var lowg float32 = 0

func Update(gameState *GameState) {
	lowg += 3.1415 / 10
	gameState.RocketPosition.Y = Height + float32(math.Sin(float64(lowg)))*Height
}
