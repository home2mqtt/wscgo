package devices

import (
	"testing"

	"gitlab.com/grill-tamasi/wscgo/tests"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

func checkDimmerPins(msg string, t *testing.T, io *tests.TestIo, on bool, pwm int) {
	if io.Values[0] != on || io.Pwm[1] != pwm {
		t.Errorf("%s ON[exp-actal]: %t - %t, PWM[exp-actal]: %d - %d\n", msg, on, io.Values[0], pwm, io.Pwm[1])
	}
}

func createDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:   0,
		PwmPin:  1,
		Speed:   10,
		OnDelay: 5,
	}
	d, _ := CreateDimmer(io, c).(*dimmer)
	return d, io
}

func createInvertedDimmerForTest() (*dimmer, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	c := &DimmerConfig{
		OnPin:    0,
		PwmPin:   1,
		Speed:    10,
		OnDelay:  5,
		Inverted: true,
	}
	d, _ := CreateDimmer(io, c).(*dimmer)
	return d, io
}

func TestDimmerInit(t *testing.T) {
	d, io := createDimmerForTest()
	d.Initialize()
	if io.Modes[0] != wiringpi.OUTPUT {
		t.Error("On Pin is not set to OUTPUT!")
	}
	if io.Modes[1] != wiringpi.PWM_OUTPUT {
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
