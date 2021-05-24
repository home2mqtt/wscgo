package devices

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

// IOutput poin configured as active output
type IOutput interface {
	Device
	Out(l gpio.Level) error
}

// OutputConfig hods the configuration for an IOutput device
type OutputConfig struct {
	Pin string `ini:"pin"`
}

type output struct {
	gpio.PinOut
}

// IInputListener function is called on changes of input
type IInputListener func(bool)

// IInput denotes a pin configued for digital input
type IInput interface {
	Device
	AddListener(IInputListener)
}

// InputConfig holds configuration for an IIinput device
type InputConfig struct {
	Pin string `ini:"pin"`
}

type input struct {
	gpio.PinIn

	listeners []IInputListener
	state     gpio.Level
}

// CreateOutput configures an IOutput device
func CreateOutput(config *OutputConfig) (IOutput, error) {
	pin := gpioreg.ByName(config.Pin)
	if pin == nil {
		return nil, invalidPinError(config.Pin)
	}
	pin = CachedPin(pin, 100)
	return &output{
		PinOut: pin,
	}, nil
}

// CreateInput configures an IInputDevice
func CreateInput(config *InputConfig) (IInput, error) {
	pin := gpioreg.ByName(config.Pin)
	if pin == nil {
		return nil, invalidPinError(config.Pin)
	}
	return &input{
		PinIn: gpioreg.ByName(config.Pin),
	}, nil
}

func (out *output) Initialize() error {
	return out.Out(gpio.Low)
}

func (*output) Tick() error {
	return nil
}

func (input *input) Initialize() error {
	err := input.In(gpio.Float, gpio.NoEdge)
	input.state = input.Read()
	return err
}

func (input *input) AddListener(listener IInputListener) {
	input.listeners = append(input.listeners, listener)
}

func (input *input) Tick() error {
	state := input.Read()
	if state != input.state {
		input.state = state
		for _, l := range input.listeners {
			l(state == gpio.High)
		}
	}
	return nil
}
