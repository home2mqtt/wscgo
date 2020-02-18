// +build windows

package wiringpi

import (
	"gitlab.com/grill-tamasi/wscgo/tests"
)

// INPUT = 0
const INPUT = 0

// OUTPUT = 1
const OUTPUT = 1

const PWM_OUTPUT = 2
const SOFT_PWM_OUTPUT = 4

type WiringPiIO struct {
	tests.TestIo
}

func Mcp23017Setup(config *Mcp23017Config) {}

func Pca9685Setup(config *Pca9685Config) {}

func (w *WiringPiIO) Setup() {
	w.TestIo = *tests.CreateTestIo(12)
}

func (*WiringPiIO) DigitalWrite(pin int, value bool) {
}

func (*WiringPiIO) DigitalRead(pin int) bool {
	return false
}

func (*WiringPiIO) PinMode(pin int, mode int) {
}

func (*WiringPiIO) PwmWrite(pin int, value int) {
}
