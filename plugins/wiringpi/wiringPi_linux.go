// +build linux

package main

// #cgo LDFLAGS: -L${SRCDIR} -L/usr/local/lib -lwiringPiPca9685 -lwiringPi
// #include<wiringPi.h>
// #include<mcp23017.h>
// #include "pca9685.h"
// #include<softPwm.h>
// extern int wiringPiDebug;
import "C"
import (
	"fmt"
	"strconv"
)

var Mcp23017rc string
var Onboardhwpwm string

// INPUT = wiringPI INPUT
const INPUT = C.INPUT

// OUTPUT = wiringPI OUTPUT
const OUTPUT = C.OUTPUT

const PWM_OUTPUT = C.PWM_OUTPUT

const onboard_pins = 64

type WiringPiIO struct {
}

func Mcp23017Setup(config *Mcp23017Config) error {
	rc := C.mcp23017Setup(C.int(config.ExpansionBase), C.int(config.Address))
	i, err := strconv.Atoi(Mcp23017rc)
	if err != nil {
		return err
	}
	if rc != (C.int)(i) {
		return fmt.Errorf("MCP23017: error  %d", rc)
	}
	return nil
}

func Pca9685Setup(config *Pca9685Config) error {
	rc := C.pca9685Setup(C.int(config.ExpansionBase), C.int(config.Address), C.float(config.Frequency))
	if rc < 0 {
		return fmt.Errorf("PCA9685: error  %d", rc)
	}
	return nil
}

func (*WiringPiIO) Setup() {
	//C.wiringPiDebug = (C.int)(1)
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

func (*WiringPiIO) PinMode(pin int, mode int) error {
	i, err := strconv.Atoi(Onboardhwpwm)
	if err != nil {
		return err
	}
	if (mode == PWM_OUTPUT) && (pin != int(i)) && (pin < onboard_pins) {
		C.softPwmCreate((C.int)(pin), 0, 1023)
	} else {
		C.pinMode((C.int)(pin), (C.int)(mode))
	}
	return nil
}

func (*WiringPiIO) PwmWrite(pin int, value int) error {
	i, err := strconv.Atoi(Onboardhwpwm)
	if err != nil {
		return err
	}
	if (pin != int(i)) && (pin < onboard_pins) {
		C.softPwmWrite((C.int)(pin), (C.int)(value))
	} else {
		C.pwmWrite((C.int)(pin), (C.int)(value))
	}
	return nil
}

type pinRange struct {
	*WiringPiIO
	start  int
	count  int
	pwmres int
}

func (pr *pinRange) PinRange() (int, int) {
	return pr.start, pr.count
}

func (pr *pinRange) PwmResolution() int {
	return pr.pwmres
}

var wiringpiio = &WiringPiIO{}

func init() {
	wiringpiio.Setup()
}
