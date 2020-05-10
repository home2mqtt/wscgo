package devices

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
)

const frequency physic.Frequency = 1000

type DimmerConfig struct {
	PwmPin     string `ini:"pwmpin"`
	OnPin      string `ini:"onpin"`
	Speed      int    `ini:"speed"`
	OnDelay    int    `ini:"ondelay"`
	Inverted   bool   `ini:"inverted"`
	Resolution int    `ini:"resolution"`
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
	BrightnessResolution() int
}

func CreateDimmer(config *DimmerConfig) (IDimmer, error) {
	var onpin gpio.PinIO
	if config.OnPin != "" {
		onpin = gpioreg.ByName(config.OnPin)
		if onpin == nil {
			return nil, invalidPinError(config.OnPin)
		}
		onpin = CachedPin(onpin, 1000)
	}
	pwmpin := gpioreg.ByName(config.PwmPin)
	if pwmpin == nil {
		return nil, invalidPinError(config.PwmPin)
	}
	pwmpin = CachedPin(pwmpin, 1000)
	return &dimmer{
		onPin:  onpin,
		pwmPin: pwmpin,
		config: config,
	}, nil
}

func (dimmer *dimmer) Initialize() error {
	dimmer.current = 0
	dimmer.target = 0
	dimmer.delaycounter = 0
	err := dimmer.pwmPin.PWM(0, frequency)
	if err != nil {
		return err
	}
	if dimmer.onPin != nil {
		return dimmer.onPin.Out(gpio.Low)
	}
	return nil
}

func (dimmer *dimmer) BrightnessResolution() int {
	return dimmer.config.Resolution
}

func (dimmer *dimmer) On() {
	dimmer.SetBrightness(dimmer.config.Resolution - 1)
}

func (dimmer *dimmer) Off() {
	dimmer.SetBrightness(0)
}

func (dimmer *dimmer) SetBrightness(target int) {
	if target > dimmer.config.Resolution-1 {
		dimmer.target = dimmer.config.Resolution - 1
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

func (dimmer *dimmer) actuate() error {
	pwmvalue := dimmer.current
	if dimmer.config.Inverted {
		pwmvalue = (dimmer.config.Resolution - 1) - pwmvalue
	}
	scaling := int(gpio.DutyMax) / dimmer.BrightnessResolution()
	err := dimmer.pwmPin.PWM(gpio.Duty(pwmvalue*scaling), frequency)
	if err != nil {
		return err
	}
	if dimmer.onPin != nil {
		l := gpio.Low
		if (dimmer.target > 0) || (dimmer.current > 0) {
			l = gpio.High
		}
		return dimmer.onPin.Out(l)
	}
	return nil
}

func (dimmer *dimmer) Tick() error {
	dimmer.adjustCurrent()
	return dimmer.actuate()
}
