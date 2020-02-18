package wiringpi

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/pin"
)

type wiringPin struct {
	ioImpl
	wpiID int
}

func New(wpiID int) gpio.PinIO {
	return &wiringPin{ioImpl: &WiringPiIO{}, wpiID: wpiID}
}

func (wp *wiringPin) String() string {
	return wp.Name()
}
func (wp *wiringPin) Halt() error {
	return nil
}

func (wp *wiringPin) Name() string {
	return fmt.Sprintf("wpi_%d", wp.wpiID)
}
func (wp *wiringPin) Number() int {
	return wp.wpiID
}
func (wp *wiringPin) Function() string {
	return "wiringPi"
}
func (wp *wiringPin) In(pull gpio.Pull, edge gpio.Edge) error {
	wp.PinMode(wp.wpiID, INPUT)
	return nil
}

func (wp *wiringPin) Read() gpio.Level {
	if wp.DigitalRead(wp.wpiID) {
		return gpio.High
	}
	return gpio.Low
}

func (wp *wiringPin) WaitForEdge(timeout time.Duration) bool {
	return false
}
func (wp *wiringPin) Pull() gpio.Pull {
	return gpio.Float
}
func (wp *wiringPin) Out(l gpio.Level) error {
	wp.PinMode(wp.wpiID, OUTPUT)
	wp.DigitalWrite(wp.wpiID, l == gpio.High)
	return nil
}
func (wp *wiringPin) PWM(duty gpio.Duty, f physic.Frequency) error {
	wp.PinMode(wp.wpiID, PWM_OUTPUT)
	wp.PwmWrite(wp.wpiID, int(duty))
	return nil
}
func (wp *wiringPin) DefaultPull() gpio.Pull {
	return gpio.Float
}

func (wp *wiringPin) Func() pin.Func {
	return gpio.IN
}
func (wp *wiringPin) SupportedFuncs() []pin.Func {
	return []pin.Func{gpio.IN, gpio.OUT, gpio.PWM}
}
func (wp *wiringPin) SetFunc(f pin.Func) error {
	return nil
}
