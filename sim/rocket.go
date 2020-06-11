package sim

import (
	"fmt"
	"math"
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
	rcsAngularMomentumChangePerTick     = 0.001                                 // kg m^2 s^-1
	ascentionFrames                     = ascentTime * 60
	// Actual physics constants
	g = 9.8
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
	frames                int
	ascending             bool
}

// ================ ROCKET STRUCT HELPERS

// CreateRocket creates and returns an instace of rocket
func CreateRocket() *Rocket {
	return &Rocket{
		Position:              Point{X: 0, Y: RocketLenght / 2},
		LiftoffTime:           time.Now(), // Simulation starts with liftoff, therefore this is appropriate
		EngineStartsRemaining: 3,          // Falcon 9 v1.1 Merlin 1D's can ignite at least 3 times https://space.stackexchange.com/questions/13953/how-do-the-falcon-9-engines-re-ignite
		fuel:                  wetMass - dryMass,
		Direction:             math.Pi / 2.0,
		ascending:             true,
	}
}

// Update the rocket to it's next physics frame
func (r *Rocket) Update() {
	r.frames++

	// Rocket orientation
	r.addRCS()
	r.updateDirection()

	// Rocket position
	r.applyGravity()
	r.addThrust()
	r.updatePosition()

	// Upkeep
	r.tickFuel()
}

// ================ ROCKET EXTERNAL FUNCTIONS

// Consts below must be pairwise coprime, they are responsible for the seamingly randomness in random seed
// Altering them means breaking current seed logic
const (
	cA = 5.0
	cB = 3.0
	cC = 2.0
)

// Ascend changes ascetion parameters randomly given seed
//
// Seed == 1 goes straight up
//
// Any other seed generates pseudorandom, coherent and repeatable behaviour for any given input
func (r *Rocket) Ascend(seed float32) {
	duration := (helpers.Sinf32(cC*seed*seed)+1.0)*ascentionFrames/5 + ascentionFrames
	if r.frames > int(duration) {
		r.SetThrust(0)
		r.ascending = false
		return
	}
	if seed == 1 {
		r.SetThrust(1)
		return
	}
	// Random thust varying from [0.8, 1.0]
	newThrust := (helpers.Sinf32(cA*seed*r.ThrustPercentage())+1.0)/20.0 + 0.9
	r.SetThrust(newThrust)
	// Target angle is varying from [67.5, 112.5]
	targetAngle := math.Pi/2.0 + helpers.Cosf32(cB*seed)*math.Pi/8.0
	if r.Direction > targetAngle {
		r.JetLeft()
	} else if r.Direction < targetAngle {
		r.JetRight()
	}
}

// IsAscending returns true whether rocket is in ascension
func (r *Rocket) IsAscending() bool {
	return r.ascending
}

// JetLeft turns on top left rcs jet (in relation to rocket's top)
func (r *Rocket) JetLeft() {
	r.AngularMomentum -= rcsAngularMomentumChangePerTick
}

// JetRight turns on top right rcs jet (in relation to rocket's top)
func (r *Rocket) JetRight() {
	r.AngularMomentum += rcsAngularMomentumChangePerTick
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

// BoundingBox return the four points that define the rocket rectangle
func (r *Rocket) BoundingBox() [4]Point {
	x := r.Position.X
	y := r.Position.Y
	a := r.Direction - math.Pi/2
	hor := float32(RocketLenght / 20.0)
	ver := float32(RocketLenght / 2.0)
	vecve := Vector{
		X: ver * helpers.Sinf32(a),
		Y: ver * helpers.Cosf32(a),
	}
	vecho := Vector{
		X: hor * helpers.Sinf32(a-math.Pi/2.0),
		Y: hor * helpers.Cosf32(a-math.Pi/2.0),
	}
	points := [4]Point{
		{
			X: x + vecve.X + vecho.X,
			Y: y + vecve.Y + vecho.Y,
		},
		{
			X: x + vecve.X - vecho.X,
			Y: y + vecve.Y - vecho.Y,
		},
		{
			X: x - vecve.X + vecho.X,
			Y: y - vecve.Y + vecho.Y,
		},
		{
			X: x - vecve.X - vecho.X,
			Y: y - vecve.Y - vecho.Y,
		},
	}
	return points
}

func (r *Rocket) Velocity() float32 {
	return float32(math.Hypot(float64(r.SpeedVector.X), float64(r.SpeedVector.Y)))
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

// Adds a little "friction" to rockets rotation, to simulate aerodinamics just a little
func (r *Rocket) updateDirection() {
	r.AngularMomentum *= 0.99
}

// Updates angle rocket points based on it's angular momentum
func (r *Rocket) addRCS() {
	r.Direction += r.AngularMomentum * physicsUpdateRate
}

// Updates rocket position based on it's speed vector
func (r *Rocket) updatePosition() {
	r.Position.X += r.SpeedVector.X * physicsUpdateRate
	r.Position.Y += r.SpeedVector.Y * physicsUpdateRate
	if r.Position.Y < 0 {
		r.Position.Y = 0
	}
}

// Applies G force on rocket (also important to remember as a
// small touch to the simulation that the gravity acceleration
// ticks down very slowly as you go up and away from earth, so
// much so that for simulation aspects, let's pretend it
// remains constant)
// It does not apply the force if rocket is on the ground
func (r *Rocket) applyGravity() {
	if DetectGroundCollision(r) >= 2 {
		r.SpeedVector.Y = 0
		return
	}
	r.SpeedVector.Y -= g * physicsUpdateRate
}

// Adds to speed vector based on current engine's thrust
func (r *Rocket) addThrust() {
	module := r.thrust * physicsUpdateRate / r.mass()
	r.SpeedVector.X += helpers.Cosf32(r.Direction) * module
	r.SpeedVector.Y += helpers.Sinf32(r.Direction) * module
}
