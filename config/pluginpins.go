package config

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/grill-tamasi/wscgo/plugins"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/pin"
)

type pluginPin struct {
	plugins.IoImpl
	wpiID   int
	pwmbits int
}

func (wp *pluginPin) String() string {
	return wp.Name()
}
func (wp *pluginPin) Halt() error {
	return nil
}

func (wp *pluginPin) Name() string {
	return fmt.Sprintf("%d", wp.wpiID)
}
func (wp *pluginPin) Number() int {
	return wp.wpiID
}
func (wp *pluginPin) Function() string {
	return "Plugin"
}
func (wp *pluginPin) In(pull gpio.Pull, edge gpio.Edge) error {
	wp.PinMode(wp.wpiID, 0)
	return nil
}

func (wp *pluginPin) Read() gpio.Level {
	if wp.DigitalRead(wp.wpiID) {
		return gpio.High
	}
	return gpio.Low
}

func (wp *pluginPin) WaitForEdge(timeout time.Duration) bool {
	return false
}
func (wp *pluginPin) Pull() gpio.Pull {
	return gpio.Float
}
func (wp *pluginPin) Out(l gpio.Level) error {
	wp.PinMode(wp.wpiID, 1)
	wp.DigitalWrite(wp.wpiID, l == gpio.High)
	return nil
}
func (wp *pluginPin) PWM(duty gpio.Duty, f physic.Frequency) error {
	if wp.pwmbits == 0 {
		return fmt.Errorf("PWM is not supported by pin %s", wp.Name())
	}
	// Scale down duty from 24 bits
	val := int(duty) >> (24 - wp.pwmbits)
	wp.PinMode(wp.wpiID, 2)
	wp.PwmWrite(wp.wpiID, val)
	log.Printf("Setting pwm value to %d\n", val)
	return nil
}
func (wp *pluginPin) DefaultPull() gpio.Pull {
	return gpio.Float
}

func (wp *pluginPin) Func() pin.Func {
	return gpio.IN
}
func (wp *pluginPin) SupportedFuncs() []pin.Func {
	return []pin.Func{gpio.IN, gpio.OUT, gpio.PWM}
}
func (wp *pluginPin) SetFunc(f pin.Func) error {
	return nil
}

type addonConfigurationPartParser struct {
	plugins.Addon
}

func (a *addonConfigurationPartParser) ParseConfiguration(cs ConfigurationSection, cc ConfigurationContext) error {
	c := a.CreateConfigStruct()
	err := cs.FillData(c)
	if err != nil {
		return err
	}

	cc.AddConfigInitializer(func() error {
		io, err := a.Configure(c)
		if err != nil {
			return err
		}
		start, count := io.PinRange()
		for i := 0; i < count; i++ {
			pin := &pluginPin{
				IoImpl: io,
				wpiID:  start + i,
			}
			gpioreg.Register(pin)
		}
		return nil
	})
	return nil
}

func loadAddon(addon plugins.Addon) error {
	return RegisterConfigurationPartParser(addon.GetIdentifier(), &addonConfigurationPartParser{
		Addon: addon,
	})
}
