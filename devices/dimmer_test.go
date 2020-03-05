package devices

import (
	"testing"

	"gitlab.com/grill-tamasi/wscgo/tests"
	"periph.io/x/periph/conn/gpio"
)

func checkDimmerPins(msg string, t *testing.T, io *tests.TestIo, on gpio.Level, pwm int) {
	if io.Pins[0].L != on || io.Pins[1].D != scale(pwm) {
		t.Errorf("%s ON[exp-actal]: %v - %v, PWM[exp-actal]: %v - %v\n", msg, on, io.Pins[0].L, scale(pwm), io.Pins[1].D)
	}
}

func createDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:      "Test_0",
		PwmPin:     "Test_1",
		Speed:      10,
		OnDelay:    5,
		Resolution: DimmerMaxValue + 1,
	}
	d, _ := CreateDimmer(c).(*dimmer)
	return d, io
}

func createInvertedDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:      "Test_0",
		PwmPin:     "Test_1",
		Speed:      10,
		OnDelay:    5,
		Inverted:   true,
		Resolution: DimmerMaxValue + 1,
	}
	d, _ := CreateDimmer(c).(*dimmer)
	return d, io
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
		expected := min(i*10, DimmerMaxValue)
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
		expected := DimmerMaxValue - min(i*10, DimmerMaxValue)
		checkDimmerPins("On ", t, io, true, expected)
		d.Tick()
	}
}

func TestDimmerOff(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()
	d.target = DimmerMaxValue
	d.current = DimmerMaxValue
	d.actuate()

	d.Off()

	for i := 0; i < 200; i++ {
		expected := max(DimmerMaxValue-i*10, 0)
		checkDimmerPins("Off ", t, io, expected > 0, expected)
		d.Tick()
	}
}

func TestDimmerOffInverted(t *testing.T) {
	d, io := createInvertedDimmerForTest()
	d.Initialize()
	d.target = DimmerMaxValue
	d.current = DimmerMaxValue
	d.actuate()

	d.Off()

	for i := 0; i < 200; i++ {
		expected := DimmerMaxValue - max(DimmerMaxValue-i*10, 0)
		checkDimmerPins("Off ", t, io, expected < DimmerMaxValue, expected)
		d.Tick()
	}
}
