package renderer

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/sim"
)

// Rocket exaust and side rcs particles

const (
	particleLifeSpan               = 5 // second
	maxNewParticlesPerFrame        = 1
	maxParticleSpeed               = 0.3
	maxParticleSpread              = 3
	maxDirectionSpread             = 10.0
	particleSize                   = 10
	particleScaleReductionPerFrame = 0.01
)

var (
	particles          []particle
	now                time.Time
	particleBaseImg, _ = ebiten.NewImage(particleSize, particleSize, ebiten.FilterDefault)
)

type particle struct {
	pos      sim.Point
	vec      sim.Vector // vector enpasulating speed and movement direction
	ang      float32    // radians
	creation time.Time
	scale    float32
}

func init() {
	particleBaseImg.Fill(color.RGBA{200, 200, 200, 255})
}

func (p *particle) draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(float64(p.scale), float64(p.scale))
	geom.Translate(float64(-particleSize*p.scale/2), float64(-particleSize*p.scale/2))
	geom.Rotate(float64(p.ang))
	geom.Translate(float64(p.pos.X), float64(p.pos.Y))
	// TODO: Add color change as particle gets old
	// colorm := {

	// }
	screen.DrawImage(particleBaseImg, &ebiten.DrawImageOptions{
		GeoM: geom,
		// ColorM: colorm,
	})
}

func (p *particle) update() {
	p.scale *= (1.0 - particleScaleReductionPerFrame)
	p.pos.X += p.vec.X
	p.pos.Y += p.vec.Y
}

func drawParticles(screen *ebiten.Image, rocket *sim.Rocket) {
	now = time.Now()
	if len(particles) < 500 {
		pos := rocketDrawPosition
		pos.X += (drawLenght / 2) * helpers.Cosf32(rocket.Direction)
		pos.Y += (drawLenght / 2) * helpers.Sinf32(rocket.Direction)
		direction := rocket.Direction
		speed := rocket.ThrustPercentage()
		volume := rocket.ThrustPercentage()
		particles = append(particles, createParticles(pos, direction, speed, volume)...)
	}
	newParticles := make([]particle, 0, 10)
	for _, p := range particles {
		if helpers.SubtractTimeInSeconds(p.creation, now) >= particleLifeSpan {
			continue
		}
		p.draw(screen)
		p.update()
		newParticles = append(newParticles, p)
	}
	particles = newParticles
}

func createParticles(pos sim.Point, direction float32, speed float32, volume float32) []particle {
	numNewParticles := int(math.Ceil(float64(float32(maxNewParticlesPerFrame) * volume)))
	newp := make([]particle, numNewParticles)
	for i := 0; i < numNewParticles; i++ {
		x := pos.X + (rand.Float32()-0.5)*maxParticleSpread
		y := pos.Y + (rand.Float32()-0.5)*maxParticleSpread
		vec := sim.Vector{
			X: helpers.Cosf32(direction) * float32(maxParticleSpeed) * speed * (rand.Float32()/maxDirectionSpread + 1.0),
			Y: helpers.Sinf32(direction) * float32(maxParticleSpeed) * speed * (rand.Float32()/maxDirectionSpread + 1.0),
		}
		p := particle{
			pos:      sim.Point{X: x, Y: y},
			vec:      vec,
			ang:      math.Pi * (rand.Float32()) / 2.0,
			creation: now,
			scale:    1,
		}
		newp = append(newp, p)
	}
	return newp
}
