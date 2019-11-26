// +build windows

package wiringpi

//import "fmt"

// LOW = 0
const LOW = 0

// HIGH = 1
const HIGH = 1

// INPUT = 0
const INPUT = 0

// OUTPUT = 1
const OUTPUT = 1

type WiringPiIO struct {
}

func Mcp23017Setup(config *Mcp23017Config) {
}

func (*WiringPiIO) DigitalWrite(pin int, value bool) {
	//fmt.Printf("DO[%d]:%d\n", pin, value)
}

func (*WiringPiIO) DigitalRead(pin int) int {
	return LOW
}

func (*WiringPiIO) PinMode(pin int, mode bool) {

}
