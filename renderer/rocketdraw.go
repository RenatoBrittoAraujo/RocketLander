package renderer

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/renatobrittoaraujo/rl/sim"
)

const (
	// minGroundDist     = 168             // This variable is an adjustment so the rocket touches the ground
	// rocketScaleAdjust = 10              // The higher this number, the less the rocket reduces in scale as it goes high
	// rocketSizeAdjust  = 0.3             // Scales rocket size to fit screen
	// maxXLag           = screenWidth / 2 // Drawn rocket position lags a little from actual position, as a visual feature
	drawLenght            = 120
	minRocketScale        = 0.3 // Percent
	maxScreenHeightRocket = 0.9 // Percent
	noScalingMaxHeight    = drawLenght * 1.5
)

var (
	rocketImage, _     = ebiten.NewImage(drawLenght/10, drawLenght, ebiten.FilterDefault)
	rocketDrawPosition sim.Point
)

func init() {
	rocketImage.Fill(color.White) // Bet you can't guess this color
}

func drawRocket(screen *ebiten.Image, rocket *sim.Rocket) {
	y := float64(rocket.Position.Y) * (float64(drawLenght) / float64(sim.RocketLenght))
	scale := rocketScale(y)

	pos := ebiten.GeoM{}
	pos.Scale(scale, scale)
	pos.Translate(-drawLenght/20, -drawLenght/2) // Adjust image positioning to center of rocket
	pos.Rotate(float64(rocket.Direction - math.Pi/2))

	posX := float64(screenWidth / 2)
	posY := float64(screenHeight)*(1.0-groundSlicePercentage) - rocketYPos(y)
	pos.Translate(posX, posY)

	rocketDrawPosition = sim.Point{X: float32(posX), Y: float32(posY)}

	screen.DrawImage(rocketImage, &ebiten.DrawImageOptions{GeoM: pos})
}

func rocketScale(h float64) float64 {
	if h <= noScalingMaxHeight {
		return 1
	}
	// These magic numbers represent a logistic function

	// These numbers are so ad hoc that i've decided to
	// leave them magic
	// But if you want to change this using the process I did,
	// find coeficients that feed your need at here
	// https://www.desmos.com/calculator
	// using vvv that as base
	return 0.2/(1+math.Pow(1.0009, h-4000)) + 0.8
}

func rocketYPos(h float64) float64 {
	if h < noScalingMaxHeight {
		return h
	}
	// These magic numbers make sure that the curve is smooth as it grows and
	// the line in the plot x=h y=screenheight passes at the point (180, 180)

	// These numbers are so ad hoc that i've decided to
	// leave them magic
	// but if you want to change this using the process I did,
	// find coeficients that feed your need at here
	// https://www.desmos.com/calculator
	// using vvv that as base
	return 540 - 130000.0/(h+181)
}
