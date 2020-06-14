package appmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/renatobrittoaraujo/rl/helpers"
	"github.com/renatobrittoaraujo/rl/input"
	"github.com/renatobrittoaraujo/rl/sim"
)

type landingLog struct {
	ID              int
	Timestamp       string
	Seed            int
	Fps             int
	Flighttime      float64 // seconds
	Score           float32
	X               float32
	Y               float32
	VerticalSpeed   float32
	HorizontalSpeed float32
	Fuel            float32
	Direction       float32
	LandingThrust   float32
}

type landingsLog struct {
	NewID          int
	AIInput        []landingLog
	HardcodedInput []landingLog
	UserInput      []landingLog
}

func logLanding(rocket *sim.Rocket, inputType int, fps int, seed int) {
	file, err := ioutil.ReadFile("logs/landing_logs.json")
	if err != nil {
		fmt.Println("Log was not found, creating file...")
		jsonFile, _ := json.Marshal(landingsLog{0, []landingLog{}, []landingLog{}, []landingLog{}})
		err = ioutil.WriteFile("logs/landing_logs.json", []byte(jsonFile), 0644)
		file, _ = ioutil.ReadFile("logs/landing_logs.json")
		if err != nil {
			fmt.Println("Could not create file 'logs/lading_logs.json', exiting logger")
			fmt.Println(err.Error())
			return
		}
	}

	var logs landingsLog
	err = json.Unmarshal(file, &logs)

	if err != nil {
		fmt.Println("Log file could not be parsed correctly, aborting")
		fmt.Println(err.Error())
		return
	}

	newid := logs.NewID
	logs.NewID++

	log := landingLog{
		ID:              newid,
		Score:           sim.LandingScore(rocket),
		X:               rocket.Position.X,
		Y:               rocket.Position.Y,
		VerticalSpeed:   rocket.SpeedVector.Y,
		HorizontalSpeed: rocket.SpeedVector.X,
		Fuel:            rocket.FuelPercentage(),
		Direction:       rocket.Direction,
		LandingThrust:   rocket.ThrustPercentage(),
		Fps:             fps,
		Seed:            seed,
		Timestamp:       time.Now().Format("2006-01-02T15:04:05.999999-07:00"),
		Flighttime:      helpers.SubtractTimeInSeconds(rocket.LiftoffTime, time.Now()),
	}

	switch inputType {
	case input.AIInput:
		logs.AIInput = append(logs.AIInput, log)
	case input.HardcodedInput:
		logs.HardcodedInput = append(logs.HardcodedInput, log)
	case input.UserInput:
		logs.UserInput = append(logs.UserInput, log)
	}

	newJSON, err := json.Marshal(logs)
	if err != nil {
		fmt.Println("New log could not be parsed correctly, aborting log")
		fmt.Println(err.Error())
		return
	}

	err = ioutil.WriteFile("logs/landing_logs.json", []byte(newJSON), 0644)
}
