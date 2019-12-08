package devices

import "gitlab.com/grill-tamasi/wscgo/wiringpi"

type IOutput interface {
	Device
	SetValue(bool)
}

type OutputConfig struct {
	Pin int `ini:"pin"`
}

type output struct {
	wiringpi.IoContext
	*OutputConfig
}

type IInputListener func(bool)

type IInput interface {
	Device
	AddListener(IInputListener)
}

type InputConfig struct {
	Pin int `ini:"pin"`
}

type input struct {
	wiringpi.IoContext
	*InputConfig

	listeners []IInputListener
	state     bool
}

func CreateOutput(io wiringpi.IoContext, config *OutputConfig) IOutput {
	return &output{
		IoContext:    io,
		OutputConfig: config,
	}
}

func CreateInput(io wiringpi.IoContext, config *InputConfig) IInput {
	return &input{
		IoContext:   io,
		InputConfig: config,
	}
}

func (out *output) Initialize() {
	out.PinMode(out.Pin, wiringpi.OUTPUT)
}

func (*output) Tick() {}

func (out *output) SetValue(value bool) {
	out.DigitalWrite(out.Pin, value)
}

func (input *input) Initialize() {
	input.PinMode(input.Pin, wiringpi.INPUT)
	input.state = input.DigitalRead(input.Pin)
}

func (input *input) AddListener(listener IInputListener) {
	input.listeners = append(input.listeners, listener)
}

func (input *input) Tick() {
	state := input.DigitalRead(input.Pin)
	if state != input.state {
		input.state = state
		for _, l := range input.listeners {
			l(state)
		}
	}
}
