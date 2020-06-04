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
	maxEngineOnTime                     = 120
	fuelComsumptionPerSecondAtMaxThrust = (wetMass - dryMass) / maxEngineOnTime // Kg
	rcsRotationTorque                   = 10.2                                  // netwon meters
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
	RotationTorque        float32
	fuel                  float32
	LiftoffTime           time.Time
	thrust                float32
	EngineStartsRemaining int
}

// ================ ROCKET STRUCT HELPERS

// CreateRocket creates and returns an instace of rocket
func CreateRocket() *Rocket {
	return &Rocket{
		LiftoffTime:           time.Now(),
		EngineStartsRemaining: 100, // Falcon 9 v1.1 Merlin 1D's can ignite at least 3 times https://space.stackexchange.com/questions/13953/how-do-the-falcon-9-engines-re-ignite
		fuel:                  wetMass - dryMass,
	}
}

// UpdateRocket the rocket to it's next physics frame
func UpdateRocket(rocket *Rocket) {
	rocket.Position.Y++
	if rocket.IsAscending() {
		// fmt.Println("ASCENT TIME")
		// Thrust max up
		//
	} else {
		// fmt.Println("LANDING TIME")
	}
	if rocket.fuel <= 0.0 {
		rocket.SetThrust(0)
	}
}

// ================ ROCKET EXTERNAL FUNCTIONS

// IsAscending returns true whether rocket is in ascension
func (g *Rocket) IsAscending() bool {
	timeDiff := helpers.SubtractTimeInSeconds(g.LiftoffTime, time.Now())
	return timeDiff < ascentTime
}

// JetLeft turns on top left rcs jet (in relation to rocket's top)
func (g *Rocket) JetLeft() {
	g.RotationTorque += rcsRotationTorque
}

// JetRight turns on top right rcs jet (in relation to rocket's top)
func (g *Rocket) JetRight() {
	g.RotationTorque -= rcsRotationTorque
}

// SetThrust sets rocket thrust to a percentage from [0.0,100.0]
func (g *Rocket) SetThrust(percentage float32) (err bool) {
	if percentage < 0 || percentage > 100 {
		panic("Input out of bounds for State.SetThrust (" + fmt.Sprintf("%0.1f", percentage) + ")")
	}
	if g.EngineStartsRemaining == 0 || g.fuel <= 0 {
		return false
	}
	if g.thrust == 0 {
		g.EngineStartsRemaining--
	}
	g.thrust = maxEngineThrust * (percentage / 100.0)
	return false
}

// FuelPercentage returns percentage from [0.0, 100.0] of fuel in rocket
func (g *Rocket) FuelPercentage() float32 {
	return g.fuel / (wetMass - dryMass)
}

// ThrustPercentage returns percentage from [0.0, 100.0] of thrust
func (g *Rocket) ThrustPercentage() float32 {
	return g.thrust / maxEngineThrust
}

// ================ ROCKET INTERNAL FUNCTIONS

// mass of rocket in kilograms
func (g *Rocket) mass() float32 {
	return g.fuel + dryMass
}

// tickFuel reduces fuel mass by current thurst amount
func (g *Rocket) tickFuel() {
	g.fuel -= g.ThrustPercentage() * fuelComsumptionPerSecondAtMaxThrust
	if g.fuel < 0 {
		g.fuel = 0
		g.EngineStartsRemaining = 0
	}
}
