package tests

import (
	"log"
)

type TestIo struct {
	Modes  []int
	Values []bool
}

func (io *TestIo) DigitalWrite(pin int, value bool) {
	if io.Values[pin] != value {
		log.Printf("Pin %d is set to %t\n", pin, value)
	}
	io.Values[pin] = value
}

func (io *TestIo) DigitalRead(pin int) bool {
	return io.Values[pin]
}

func (io *TestIo) PinMode(pin int, mode int) {
	log.Printf("Mode of pin %d is set to %d\n", pin, mode)
	io.Modes[pin] = mode
}

func CreateTestIo(pins int) *TestIo {
	return &TestIo{
		Modes:  make([]int, pins),
		Values: make([]bool, pins),
	}
}
