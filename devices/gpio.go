package devices

type IOutput interface {
	Device
	SetValue(bool)
}

type OutputConfig struct {
	Pin int
}

type output struct {
	IoContext
	*OutputConfig
}

func CreateOutput(io IoContext, config *OutputConfig) IOutput {
	return &output{
		IoContext:    io,
		OutputConfig: config,
	}
}

func (out *output) Initialize() {
	out.PinMode(out.Pin, true)
}

func (*output) Tick() {}

func (out *output) SetValue(value bool) {
	out.DigitalWrite(out.Pin, value)
}
