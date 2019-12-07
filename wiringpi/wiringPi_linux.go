// +build linux

package wiringpi

// #cgo LDFLAGS: -lwiringPi
// #include<wiringPi.h>
// #include<mcp23017.h>
// #ifdef PI_MODEL_BPR
// int setuprc = 0;
// #else
// int setuprc = 1;
// #endif
import "C"
import "log"

// INPUT = wiringPI INPUT
const INPUT = C.INPUT

// OUTPUT = wiringPI OUTPUT
const OUTPUT = C.OUTPUT

const PWM_OUTPUT = C.PWM_OUTPUT

type WiringPiIO struct {
}

func Mcp23017Setup(config *Mcp23017Config) {
	rc := C.mcp23017Setup(C.int(config.ExpansionBase), C.int(config.Address))
	if rc != C.setuprc {
		log.Fatal("MCP23017 error: ", rc)
	}
	C.wiringPiSetup()
}

func (*WiringPiIO) DigitalWrite(pin int, value bool) {
	v := C.LOW
	if value {
		v = C.HIGH
	}
	C.digitalWrite((C.int)(pin), (C.int)(v))
}

func (*WiringPiIO) DigitalRead(pin int) bool {
	return C.HIGH == C.digitalRead((C.int)(pin))
}

func (*WiringPiIO) PinMode(pin int, mode int) {
	C.pinMode((C.int)(pin), (C.int)(mode))
}

func (*WiringPiIO) PwmWrite(pin int, value int) {
	C.pwmWrite((C.int)(pin), (C.int)(value))
}
