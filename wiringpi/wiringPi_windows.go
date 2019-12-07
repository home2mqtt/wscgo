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

type WiringPiIO struct {
	tests.TestIo
}

func Mcp23017Setup(config *Mcp23017Config) {}
