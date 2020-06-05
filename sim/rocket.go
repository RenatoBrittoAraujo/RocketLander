package sim

import (
	"fmt"
	"time"

	"github.com/renatobrittoaraujo/rl/helpers"
)

// Mass data from: https://sma.nasa.gov/LaunchVehicle/assets/spacex-falcon-9-data-sheet.pdf
const (
	// Data from Falcon 9 v1.1
	RocketLenght    = 70      // meters
	maxEngineThrust = 5885000 // newtons
	dryMass         = 28000   // kilograms
	wetMass         = 439000
	// Constants related purely with simulation
	ascentTime                          = 5 // seconds
	maxEngineOnTime                     = 100
	physicsUpdateRate                   = 1.0 / 60.0                            // per second
	fuelComsumptionPerSecondAtMaxThrust = (wetMass - dryMass) / maxEngineOnTime // kg
	rcsAngularMomentumChangePerTick     = 0.0001                                // kg m^2 s^-1
	rotation
	// Actual physics constants
	g  = 9.8
	pi = 3.1415926
)

// Rocket holds all relevant simulation data
//
// speed given by vector of meters per second
//
// direction given in radians (0 is vertical up)
//
// thrust given in newtons
//
// fuel given in kilograms
type Rocket struct {
	Position              Point
	Direction             float32
	SpeedVector           Vector
	AngularMomentum       float32
	LiftoffTime           time.Time
	EngineStartsRemaining int
	fuel                  float32
	thrust                float32
	controlsFree          bool
}

// ================ ROCKET STRUCT HELPERS

// CreateRocket creates and returns an instace of rocket
func CreateRocket() *Rocket {
	return &Rocket{
		LiftoffTime:           time.Now(),
		EngineStartsRemaining: 3, // Falcon 9 v1.1 Merlin 1D's can ignite at least 3 times https://space.stackexchange.com/questions/13953/how-do-the-falcon-9-engines-re-ignite
		fuel:                  wetMass - dryMass,
		controlsFree:          false,
		Direction:             pi / 2.0,
	}
}

// Update the rocket to it's next physics frame
func (r *Rocket) Update() {
	// r.Position.Y++

	// Rocket orientation
	r.addRCS()
	r.updateDirection()

	// Rocket position
	r.applyGravity()
	r.addThrust()
	r.updatePosition()

	// Upkeep
	r.tickFuel()

	fmt.Println(r.SpeedVector)
}

func (r *Rocket) updateDirection() {
	r.AngularMomentum *= 0.99
}

func (r *Rocket) addRCS() {
	r.Direction += r.AngularMomentum * physicsUpdateRate
}

func (r *Rocket) updatePosition() {
	r.Position.X += r.SpeedVector.X * physicsUpdateRate
	r.Position.Y += r.SpeedVector.Y * physicsUpdateRate
	if r.Position.Y < 0 {
		r.Position.Y = 0
	}
}

func (r *Rocket) applyGravity() {
	if r.Position.Y <= 0 {
		return
	}
	r.SpeedVector.Y -= g * physicsUpdateRate
}

func (r *Rocket) addThrust() {
	fmt.Println("ANG:", r.Direction*180/pi)
	module := r.thrust * physicsUpdateRate / r.mass()
	fmt.Println("THRUST MODULE:", module)
	fmt.Println("CHANGE X:", helpers.Cosf32(r.Direction)*module)
	fmt.Println("CHANGE Y:", helpers.Sinf32(r.Direction)*module)
	r.SpeedVector.X += helpers.Cosf32(r.Direction) * module
	r.SpeedVector.Y += helpers.Sinf32(r.Direction) * module
}

// ================ ROCKET EXTERNAL FUNCTIONS

// Ascend changes ascetion parameters randomly given seed
//
// Seed == 1 goes straight up
//
// Any other seed generates random (coherent) behaviour
func (r *Rocket) Ascend(seed float64) {
	r.SetThrust(1)
}

// IsAscending returns true whether rocket is in ascension
func (r *Rocket) IsAscending() bool {
	if r.controlsFree {
		return false
	}
	timeDiff := helpers.SubtractTimeInSeconds(r.LiftoffTime, time.Now())
	r.controlsFree = timeDiff >= ascentTime
	return !r.controlsFree
}

// JetLeft turns on top left rcs jet (in relation to rocket's top)
func (r *Rocket) JetLeft() {
	r.AngularMomentum += rcsAngularMomentumChangePerTick
}

// JetRight turns on top right rcs jet (in relation to rocket's top)
func (r *Rocket) JetRight() {
	r.AngularMomentum -= rcsAngularMomentumChangePerTick
}

// SetThrust sets rocket thrust to a percentage from [0.0,1.0]
func (r *Rocket) SetThrust(percentage float32) (err bool) {
	if percentage < 0 || percentage > 1 {
		panic("Input out of bounds for State.SetThrust (" + fmt.Sprintf("%0.1f", percentage) + ")")
	}
	if (r.EngineStartsRemaining == 0 && r.thrust <= 0) || r.fuel <= 0 {
		return false
	}
	if r.thrust == 0 && percentage > 0 {
		r.EngineStartsRemaining--
	}
	r.thrust = maxEngineThrust * percentage
	return false
}

// FuelPercentage returns percentage from [0.0, 1.0] of fuel in rocket
func (r *Rocket) FuelPercentage() float32 {
	return r.fuel / (wetMass - dryMass)
}

// ThrustPercentage returns percentage from [0.0, 1.0] of thrust
func (r *Rocket) ThrustPercentage() float32 {
	return r.thrust / maxEngineThrust
}

// ================ ROCKET INTERNAL FUNCTIONS

// mass of rocket in kilograms
func (r *Rocket) mass() float32 {
	return r.fuel + dryMass
}

// tickFuel reduces fuel mass by current thurst amount
func (r *Rocket) tickFuel() {
	if r.fuel <= 0 {
		r.fuel = 0
		return
	}
	r.fuel -= r.ThrustPercentage() * fuelComsumptionPerSecondAtMaxThrust * physicsUpdateRate
	if r.fuel <= 0 {
		r.fuel = 0
		r.EngineStartsRemaining = 0
		r.thrust = 0
	}
}
