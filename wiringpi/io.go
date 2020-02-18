package wiringpi

import (
	"periph.io/x/periph/conn/gpio"
)

// LOW = wiringPI LOW
const LOW = false

// HIGH = wiringPI HIGH
const HIGH = true

type IoContext interface {
	Setup()

	GetPin(pin int) gpio.PinIO
}

type ioImpl interface {
	DigitalWrite(pin int, value bool)
	DigitalRead(pin int) bool
	PinMode(pin int, mode int)
	PwmWrite(pin int, value int)
}
