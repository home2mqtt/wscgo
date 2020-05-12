package devices

import (
	"log"
	"testing"

	"github.com/balazsgrill/wscgo/tests"
	"periph.io/x/periph/conn/gpio"
)

func checkDimmerPins(msg string, t *testing.T, io *tests.TestIo, on gpio.Level, pwm gpio.Duty) {
	if io.Pins[0].L != on || io.Pins[1].D != pwm {
		t.Errorf("%s ON[exp-actal]: %v - %v, PWM[exp-actal]: %v - %v\n", msg, on, io.Pins[0].L, pwm, io.Pins[1].D)
	}
}

func createDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:   "Test_0",
		PwmPin:  "Test_1",
		Speed:   8192,
		OnDelay: 5,
	}
	id, _ := CreateDimmer(c)
	d, _ := id.(*dimmer)
	return d, io
}

func createInvertedDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:    "Test_0",
		PwmPin:   "Test_1",
		Speed:    8192,
		OnDelay:  5,
		Inverted: true,
	}
	id, _ := CreateDimmer(c)
	d, _ := id.(*dimmer)
	return d, io
}

func TestScale(t *testing.T) {
	value := gpio.DutyMax
	log.Printf("%d", value)
	log.Printf("%d", value>>8)
}

func TestDimmerInit(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()
	if io.Pins[0].L != gpio.Low {
		t.Error("On Pin is not set to OUTPUT!")
	}
	if io.Pins[1].F != frequency {
		t.Error("On Pin is not set to PWM!")
	}
}

func TestDimmerDelay(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()

	d.On()
	for i := 0; i < 5; i++ {
		d.Tick()
		checkDimmerPins("Dimmer ", t, io, true, 0)
	}

}

func TestDimmerOn(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()

	d.On()
	for i := 0; i < 5; i++ {
		d.Tick()
	}

	for i := 0; i < 200; i++ {
		expected := min(gpio.Duty(i)*gpio.Duty(d.config.Speed), gpio.DutyMax)
		checkDimmerPins("On ", t, io, true, expected)
		d.Tick()
	}
}

func TestDimmerOnInverted(t *testing.T) {
	d, io := createInvertedDimmerForTest()
	d.Initialize()

	d.On()
	for i := 0; i < 5; i++ {
		d.Tick()
	}

	for i := 0; i < 200; i++ {
		expected := gpio.DutyMax - min(gpio.Duty(i)*gpio.Duty(d.config.Speed), gpio.DutyMax)
		checkDimmerPins("On ", t, io, true, expected)
		d.Tick()
	}
}

func TestDimmerOff(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()
	d.target = gpio.DutyMax
	d.current = gpio.DutyMax
	d.actuate()

	d.Off()

	for i := 0; i < 2000; i++ {
		expected := max(gpio.DutyMax-gpio.Duty(i)*gpio.Duty(d.config.Speed), 0)
		checkDimmerPins("Off ", t, io, expected > 0, expected)
		d.Tick()
	}
}

func TestDimmerOffInverted(t *testing.T) {
	d, io := createInvertedDimmerForTest()
	d.Initialize()
	d.target = gpio.DutyMax
	d.current = gpio.DutyMax
	d.actuate()

	d.Off()

	for i := 0; i < 2000; i++ {
		expected := gpio.DutyMax - max(gpio.DutyMax-gpio.Duty(i)*gpio.Duty(d.config.Speed), 0)
		checkDimmerPins("Off ", t, io, expected < gpio.DutyMax, expected)
		d.Tick()
	}
}
