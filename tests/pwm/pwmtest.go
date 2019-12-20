package main

import (
	"math"
	"time"

	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

func main() {
	io := &wiringpi.WiringPiIO{}
	io.Setup()
	angle := 0
	pin := 6

	io.PinMode(pin, wiringpi.SOFT_PWM_OUTPUT)

	controlTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		println("Timer started")
		for range controlTicker.C {
			rad := math.Pi * float64(angle) / (180)
			v := 511 + int(math.Round(math.Sin(rad)*512))
			io.PwmWrite(pin, v)
			angle++
		}
	}()

	select {}
}
