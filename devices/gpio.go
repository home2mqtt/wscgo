package devices

import "gitlab.com/grill-tamasi/wscgo/wiringpi"

type IOutput interface {
	Device
	SetValue(bool)
}

type OutputConfig struct {
	Pin int
}

type output struct {
	wiringpi.IoContext
	*OutputConfig
}

func CreateOutput(io wiringpi.IoContext, config *OutputConfig) IOutput {
	return &output{
		IoContext:    io,
		OutputConfig: config,
	}
}

func (out *output) Initialize() {
	out.PinMode(out.Pin, wiringpi.OUTPUT)
}

func (*output) Tick() {}

func (out *output) SetValue(value bool) {
	out.DigitalWrite(out.Pin, value)
}
