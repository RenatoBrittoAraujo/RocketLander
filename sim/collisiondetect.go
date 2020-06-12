package sim

import "math"

const (
	// MaxLandingVelocity keeps track of maximum landing speed (vertical + horizontal)
	MaxLandingVelocity = 20 // m/s
	// MaxAngleDeviation keeps track of maximum angle rocket can make to a perpendicular line to ground for sucessful landing
	MaxAngleDeviation = 0.872 // Radians
)

// DetectGroundCollision returns true if ground collision happend
func DetectGroundCollision(r *Rocket) (collision int) {
	points := r.BoundingBox()
	for _, v := range points {
		if v.Y <= 0 {
			collision++
		}
	}
	return
}

// LandingScore is a score of a particular landing
//
// ret >= 0 means a successful landing and < 0 unsuccessful landing
func LandingScore(r *Rocket) float32 {
	score := 1.0
	speed := r.Velocity()
	// Logistical function, below 20 speed it gives positive and negative for everything else
	// These numbers are magical and based on MaxLandingVelocity = 20
	// If you wish to change them, use these coeficients below here:
	// https://www.desmos.com/calculator
	score += 10.0/(1+math.Exp(0.6*(float64(speed)-MaxLandingVelocity*0.8))) - 5.0
	angleFromUpright := math.Abs(math.Pi/2 - float64(r.Direction))
	score += 10.0/(1+math.Exp(0.9*(angleFromUpright-MaxAngleDeviation*0.8))) - 5.0
	return float32(score)
}
