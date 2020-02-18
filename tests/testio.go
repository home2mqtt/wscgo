package tests

import (
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpiotest"
)

type TestIo struct {
	Pins []*gpiotest.Pin
}

func (io *TestIo) Setup() {

}

func (io *TestIo) GetPin(pin int) gpio.PinIO {
	return io.Pins[pin]
}

func CreateTestIo(numofpins int) *TestIo {
	pins := make([]*gpiotest.Pin, numofpins)
	for i := range pins {
		pins[i] = &gpiotest.Pin{
			N:   fmt.Sprintf("Test_%d", i),
			Num: i,
		}
	}
	return &TestIo{
		Pins: pins,
	}
}
