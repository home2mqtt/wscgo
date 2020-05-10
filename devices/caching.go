package devices

import (
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
)

type cachedGpio struct {
	pin       gpio.PinIO
	reinforce int

	count int
	pwm   bool
	level gpio.Level
	duty  gpio.Duty
}

func CachedPin(pin gpio.PinIO, reinforce int) gpio.PinIO {
	return &cachedGpio{
		pin:       pin,
		reinforce: reinforce,

		count: 0,
	}
}

func (cp *cachedGpio) String() string {
	return cp.pin.String()
}
func (cp *cachedGpio) Halt() error {
	return cp.pin.Halt()
}

func (cp *cachedGpio) Name() string {
	return cp.pin.Name()
}
func (cp *cachedGpio) Number() int {
	return cp.pin.Number()
}
func (cp *cachedGpio) Function() string {
	return cp.pin.Function()
}
func (cp *cachedGpio) In(pull gpio.Pull, edge gpio.Edge) error {
	return cp.pin.In(pull, edge)
}

func (cp *cachedGpio) Read() gpio.Level {
	return cp.pin.Read()
}

func (cp *cachedGpio) WaitForEdge(timeout time.Duration) bool {
	return cp.pin.WaitForEdge(timeout)
}
func (cp *cachedGpio) Pull() gpio.Pull {
	return cp.pin.Pull()
}
func (cp *cachedGpio) Out(l gpio.Level) error {
	if cp.count == 0 || cp.pwm || l != cp.level {
		err := cp.pin.Out(l)
		if err == nil {
			cp.pwm = false
			cp.count = cp.reinforce
			cp.level = l
		}
		return err
	} else {
		if cp.count > 0 {
			cp.count--
		}
	}
	return nil
}
func (cp *cachedGpio) PWM(duty gpio.Duty, f physic.Frequency) error {
	if cp.count == 0 || !cp.pwm || cp.duty != duty {
		err := cp.pin.PWM(duty, f)
		if err == nil {
			cp.pwm = true
			cp.count = cp.reinforce
			cp.duty = duty
		}
		return err
	} else {
		if cp.count > 0 {
			cp.count--
		}
	}
	return nil
}
func (cp *cachedGpio) DefaultPull() gpio.Pull {
	return cp.pin.DefaultPull()
}
