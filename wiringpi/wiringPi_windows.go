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
