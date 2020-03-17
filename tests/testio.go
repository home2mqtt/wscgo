package tests

import (
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
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
		pn := fmt.Sprintf("Test_%d", i)
		gpioreg.Unregister(pn)
		pins[i] = &gpiotest.Pin{
			N:   pn,
			Num: i,
		}
		gpioreg.Register(pins[i])
	}
	return &TestIo{
		Pins: pins,
	}
}
