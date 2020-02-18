package devices

import (
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
)

const DimmerResolution int = 1 << (24 - 16)
const DimmerMaxValue int = DimmerResolution - 1

const frequency physic.Frequency = 1000

type DimmerConfig struct {
	PwmPin   int  `ini:"pwmpin"`
	OnPin    int  `ini:"onpin"`
	Speed    int  `ini:"speed"`
	OnDelay  int  `ini:"ondelay"`
	Inverted bool `ini:"inverted"`
}

type dimmer struct {
	onPin        gpio.PinOut
	pwmPin       gpio.PinOut
	config       *DimmerConfig
	current      int
	target       int
	delaycounter int
}

type IDimmer interface {
	Device
	On()
	Off()
	SetBrightness(value int)
}

func CreateDimmer(io wiringpi.IoContext, config *DimmerConfig) IDimmer {
	return &dimmer{
		onPin:  io.GetPin(config.OnPin),
		pwmPin: io.GetPin(config.PwmPin),
		config: config,
	}
}

func (dimmer *dimmer) Initialize() {
	dimmer.current = 0
	dimmer.target = 0
	dimmer.delaycounter = 0
	dimmer.pwmPin.PWM(0, frequency)
	if dimmer.onPin != nil {
		dimmer.onPin.Out(gpio.Low)
	}
}

func (dimmer *dimmer) On() {
	dimmer.SetBrightness(DimmerMaxValue)
}

func (dimmer *dimmer) Off() {
	dimmer.SetBrightness(0)
}

func (dimmer *dimmer) SetBrightness(target int) {
	if target > DimmerMaxValue {
		dimmer.target = DimmerMaxValue
	} else {
		if target < 0 {
			dimmer.target = 0
		} else {
			dimmer.target = target
		}
	}

	if (dimmer.target != 0) && (dimmer.current == 0) {
		dimmer.delaycounter = dimmer.config.OnDelay
	}
}

func min(v1 int, v2 int) int {
	if v1 > v2 {
		return v2
	}
	return v1
}

func max(v1 int, v2 int) int {
	if v1 < v2 {
		return v2
	}
	return v1
}

func (dimmer *dimmer) adjustCurrent() {
	if dimmer.delaycounter > 0 {
		dimmer.delaycounter--
		return
	}
	if dimmer.target > dimmer.current {
		dimmer.current = min(dimmer.target, dimmer.current+dimmer.config.Speed)
		return
	}
	if dimmer.target < dimmer.current {
		dimmer.current = max(dimmer.target, dimmer.current-dimmer.config.Speed)
		return
	}
}

func scale(brightness int) gpio.Duty {
	return gpio.Duty(int32(brightness) * int32(DimmerResolution))
}

func (dimmer *dimmer) actuate() {
	pwmvalue := dimmer.current
	if dimmer.config.Inverted {
		pwmvalue = DimmerMaxValue - pwmvalue
	}
	dimmer.pwmPin.PWM(scale(pwmvalue), frequency)
	if dimmer.onPin != nil {
		l := gpio.Low
		if (dimmer.target > 0) || (dimmer.current > 0) {
			l = gpio.High
		}
		dimmer.onPin.Out(l)
	}
}

func (dimmer *dimmer) Tick() {
	dimmer.adjustCurrent()
	dimmer.actuate()
}
