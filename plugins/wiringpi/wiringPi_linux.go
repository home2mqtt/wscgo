// +build linux

package main

// #cgo LDFLAGS: -L${SRCDIR} -L/usr/local/lib -lwiringPiPca9685 -lwiringPi
// #include<wiringPi.h>
// #include<mcp23017.h>
// #include "pca9685.h"
// #include<softPwm.h>
// #ifdef PI_MODEL_BPR
// int setuprc = 0;
// int onboard_hw_pwm = -1;
// #else
// int setuprc = 1;
// int onboard_hw_pwm = 1;
// #endif
// extern int wiringPiDebug;
import "C"
import (
	"fmt"
	"log"
)

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
	if rc != C.setuprc {
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

func (*WiringPiIO) PinMode(pin int, mode int) {
	if (mode == PWM_OUTPUT) && ((C.int)(pin) != C.onboard_hw_pwm) && (pin < onboard_pins) {
		C.softPwmCreate((C.int)(pin), 0, 1023)
	} else {
		C.pinMode((C.int)(pin), (C.int)(mode))
	}
}

func (*WiringPiIO) PwmWrite(pin int, value int) {
	if ((C.int)(pin) != C.onboard_hw_pwm) && (pin < onboard_pins) {
		C.softPwmWrite((C.int)(pin), (C.int)(value))
	} else {
		log.Printf("pwmWrite(%d, %d)\n", pin, value)
		C.pwmWrite((C.int)(pin), (C.int)(value))
	}
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
