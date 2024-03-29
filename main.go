package main

import (
	"log"
	"play-go/controls"
	"play-go/fixture"
	"time"

	"github.com/akualab/dmx"
)

func main() {
	connection, err := dmx.NewDMXConnection("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}

	freedomPars := fixture.NewFixtureList()
	for _, address := range []int{65, 81, 97, 113} {
		freedomPars.AddFixture(fixture.NewFreedomPar(), address)
	}
	tomshine := fixture.NewFixtureList()
	for _, address := range []int{1, 17, 33, 49} {
		tomshine.AddFixture(fixture.NewTomeshine(), address)
	}

	universe := fixture.NewFixtureList()
	universe.AddFixtureList(freedomPars)
	universe.AddFixtureList(tomshine)

	universe.SetAll(0)
	universe.SetValue("dimmer", 0)

	soft_white := controls.AllColors["soft_white"]
	soft_white.Values().ApplyTo(universe)

	surface := NewControls()

	dimmerDial := surface["dimmer"]
	controls.LinkFixtureChannel(universe, "dimmer", dimmerDial.Channel)

	tiltDial := surface["tilt"]
	controls.LinkFixtureChannel(universe, "tilt", tiltDial.Channel)

	// Repeat(8*time.Second, GetToggleFunc(dimmerDial, 6*time.Second))

	Render(universe, connection)

	// universe.SetValue("dimmer", 255)

	StartServer(surface)

	time.Sleep(100 * time.Second)
}

func Repeat(duration time.Duration, f func()) *time.Ticker {
	f()
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}

func GetToggleFunc(dial *controls.Dial, duration time.Duration) func() {
	return func() {
		endValue := dial.Opposite()
		log.Printf("triggering ease from %d to %d", dial.Value, endValue)
		Ease(dial, duration/2, endValue)
	}
}

func Ease(dial *controls.Dial, duration time.Duration, endValue byte) {
	startValue := float64(dial.Value)
	startTime := time.Now()
	ticker := time.NewTicker(REFRESH).C
	go func() {
		for t := range ticker {
			if t.After(startTime.Add(duration)) {
				dial.SetValue(endValue)
				return
			}
			elapsed := t.Sub(startTime)
			step := float64(elapsed) / float64(duration)
			// mult := ease.InOutSine(step)
			mult := step
			value := startValue + (float64(endValue)-startValue)*mult

			log.Printf("elapsed: %f step: %f mult: %f value: %f",
				float64(elapsed), step, mult, value,
			)

			dial.SetValue(byte(value))
		}
	}()
}

func NewControls() controls.DialMap {
	m := make(controls.DialMap)
	m["dimmer"] = controls.NewDial()
	m["tilt"] = controls.NewDial()
	return m
}

const REFRESH = 11 * time.Millisecond

func Render(f fixture.FixtureList, connection *dmx.DMX) {
	ticker := time.NewTicker(REFRESH)
	go func() {
		for range ticker.C {
			f.Render(*connection)
		}
	}()
}
