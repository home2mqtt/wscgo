// +build windows

package wiringpi

import "fmt"

// INPUT = 0
const INPUT = 0

// OUTPUT = 1
const OUTPUT = 1

type WiringPiIO struct {
}

func Mcp23017Setup(config *Mcp23017Config) {
}

func (*WiringPiIO) DigitalWrite(pin int, value bool) {
	fmt.Printf("DO[%d]:%t\n", pin, value)
}

func (*WiringPiIO) DigitalRead(pin int) bool {
	return LOW
}

func (*WiringPiIO) PinMode(pin int, mode int) {

}
