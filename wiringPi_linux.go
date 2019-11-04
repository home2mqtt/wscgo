// +build linux

package main

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

// LOW = wiringPI LOW
const LOW = C.LOW

// HIGH = wiringPI HIGH
const HIGH = C.HIGH

// INPUT = wiringPI INPUT
const INPUT = C.INPUT

// OUTPUT = wiringPI OUTPUT
const OUTPUT = C.OUTPUT

type wiringPiIO struct {
}

func mcp23017Setup(expansionBase int, address int) {
	rc := C.mcp23017Setup(C.int(expansionBase), C.int(address))
	if rc != C.setuprc {
		log.Fatal("MCP23017 error: ", rc)
	}
	C.wiringPiSetup()
}

func (*wiringPiIO) digitalWrite(pin int, value int) {
	C.digitalWrite((C.int)(pin), (C.int)(value))
}

func (*wiringPiIO) digitalRead(pin int) int {
	return int(C.digitalRead((C.int)(pin)))
}

func (*wiringPiIO) pinMode(pin int, mode int) {
	C.pinMode((C.int)(pin), (C.int)(mode))
}
