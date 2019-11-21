package tests

import (
	"log"

	"gitlab.com/grill-tamasi/wscgo/devices"
)

type testIo struct {
	modes  []bool
	values []bool
}

func (io *testIo) DigitalWrite(pin int, value bool) {
	if io.values[pin] != value {
		log.Printf("Pin %d is set to %t\n", pin, value)
	}
	io.values[pin] = value
}

func (io *testIo) DigitalRead(pin int) bool {
	return io.values[pin]
}

func (io *testIo) PinMode(pin int, mode bool) {
	log.Printf("Mode of pin %d is set to %t\n", pin, mode)
	io.modes[pin] = mode
}

func CreateTestIo(pins int) devices.IoContext {
	return &testIo{
		modes:  make([]bool, pins),
		values: make([]bool, pins),
	}
}
