package devices

import (
	"testing"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpiotest"
	"periph.io/x/conn/v3/physic"
)

type callCountPin struct {
	*gpiotest.Pin

	OutCalls int
	PwmCalls int
}

func (p *callCountPin) Out(l gpio.Level) error {
	p.OutCalls++
	return p.Pin.Out(l)
}

func (p *callCountPin) PWM(duty gpio.Duty, f physic.Frequency) error {
	p.PwmCalls++
	return p.Pin.PWM(duty, f)
}

func TestCached_Out(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	put.Out(gpio.High)
	put.Out(gpio.High)
	if pin.L != gpio.High {
		t.Error("Pin value is invalid")
	}
	if pin.OutCalls != 1 {
		t.Errorf("Out is called %d times", pin.OutCalls)
	}
}

func TestCached_Out2(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	put.Out(gpio.High)
	put.Out(gpio.Low)
	if pin.L != gpio.Low {
		t.Error("Pin value is invalid")
	}
	if pin.OutCalls != 2 {
		t.Errorf("Out is called %d times", pin.OutCalls)
	}
}

func TestCached_Out3(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	for i := 0; i < 23; i++ {
		put.Out(gpio.High)
	}
	if pin.L != gpio.High {
		t.Error("Pin value is invalid")
	}
	if pin.OutCalls != 3 {
		t.Errorf("Out is called %d times", pin.OutCalls)
	}
}

func TestCached_Pwm(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	put.PWM(100, 100)
	put.PWM(100, 100)
	if pin.D != 100 {
		t.Error("PWM value is invalid")
	}
	if pin.PwmCalls != 1 {
		t.Errorf("PWM is called %d times", pin.PwmCalls)
	}
}

func TestCached_Pwm2(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	put.PWM(200, 100)
	put.PWM(100, 100)
	if pin.D != 100 {
		t.Error("PWM value is invalid")
	}
	if pin.PwmCalls != 2 {
		t.Errorf("PWM is called %d times", pin.PwmCalls)
	}
}

func TestCached_Pwm3(t *testing.T) {
	pin := &callCountPin{
		Pin: &gpiotest.Pin{},
	}
	put := CachedPin(pin, 10)

	for i := 0; i < 23; i++ {
		put.PWM(100, 100)
	}
	if pin.D != 100 {
		t.Error("PWM value is invalid")
	}
	if pin.PwmCalls != 3 {
		t.Errorf("PWM is called %d times", pin.PwmCalls)
	}
}
