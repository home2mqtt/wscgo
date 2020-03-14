package devices

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type IOutput interface {
	Device
	Out(l gpio.Level) error
}

type OutputConfig struct {
	Pin string `ini:"pin"`
}

type output struct {
	gpio.PinOut
}

type IInputListener func(bool)

type IInput interface {
	Device
	AddListener(IInputListener)
}

type InputConfig struct {
	Pin string `ini:"pin"`
}

type input struct {
	gpio.PinIn

	listeners []IInputListener
	state     gpio.Level
}

func CreateOutput(config *OutputConfig) (IOutput, error) {
	pin := gpioreg.ByName(config.Pin)
	if pin == nil {
		return nil, invalidPinError(config.Pin)
	}
	return &output{
		PinOut: pin,
	}, nil
}

func CreateInput(config *InputConfig) (IInput, error) {
	pin := gpioreg.ByName(config.Pin)
	if pin == nil {
		return nil, invalidPinError(config.Pin)
	}
	return &input{
		PinIn: gpioreg.ByName(config.Pin),
	}, nil
}

func (out *output) Initialize() {
	out.Out(gpio.Low)
}

func (*output) Tick() {}

func (input *input) Initialize() {
	input.In(gpio.Float, gpio.NoEdge)
	input.state = input.Read()
}

func (input *input) AddListener(listener IInputListener) {
	input.listeners = append(input.listeners, listener)
}

func (input *input) Tick() {
	state := input.Read()
	if state != input.state {
		input.state = state
		for _, l := range input.listeners {
			l(state == gpio.High)
		}
	}
}
