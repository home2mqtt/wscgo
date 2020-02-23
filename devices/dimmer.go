package devices

import "gitlab.com/grill-tamasi/wscgo/wiringpi"

type DimmerConfig struct {
	PwmPin     int  `ini:"pwmpin"`
	OnPin      int  `ini:"onpin"`
	Speed      int  `ini:"speed"`
	OnDelay    int  `ini:"ondelay"`
	Inverted   bool `ini:"inverted"`
	Resolution int  `ini:"resolution"`
}

type dimmer struct {
	wiringpi.IoContext
	*DimmerConfig
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

func CreateDimmer(io wiringpi.IoContext, config *DimmerConfig) IDimmer {
	return &dimmer{
		IoContext:    io,
		DimmerConfig: config,
	}
}

func (dimmer *dimmer) Initialize() {
	dimmer.current = 0
	dimmer.target = 0
	dimmer.delaycounter = 0
	dimmer.PinMode(dimmer.PwmPin, wiringpi.PWM_OUTPUT)
	if dimmer.OnPin >= 0 {
		dimmer.PinMode(dimmer.OnPin, wiringpi.OUTPUT)
	}
}

func (dimmer *dimmer) BrightnessResolution() int {
	return dimmer.Resolution
}

func (dimmer *dimmer) On() {
	dimmer.SetBrightness(dimmer.Resolution - 1)
}

func (dimmer *dimmer) Off() {
	dimmer.SetBrightness(0)
}

func (dimmer *dimmer) SetBrightness(target int) {
	if target > dimmer.Resolution-1 {
		dimmer.target = dimmer.Resolution - 1
	} else {
		if target < 0 {
			dimmer.target = 0
		} else {
			dimmer.target = target
		}
	}

	if (dimmer.target != 0) && (dimmer.current == 0) {
		dimmer.delaycounter = dimmer.OnDelay
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
		dimmer.current = min(dimmer.target, dimmer.current+dimmer.Speed)
		return
	}
	if dimmer.target < dimmer.current {
		dimmer.current = max(dimmer.target, dimmer.current-dimmer.Speed)
		return
	}
}

func (dimmer *dimmer) actuate() {
	pwmvalue := dimmer.current
	if dimmer.Inverted {
		pwmvalue = (dimmer.Resolution - 1) - pwmvalue
	}
	dimmer.PwmWrite(dimmer.PwmPin, pwmvalue)
	if dimmer.OnPin >= 0 {
		dimmer.DigitalWrite(dimmer.OnPin, (dimmer.target > 0) || (dimmer.current > 0))
	}
}

func (dimmer *dimmer) Tick() {
	dimmer.adjustCurrent()
	dimmer.actuate()
}
