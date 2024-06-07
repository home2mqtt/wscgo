package integration

import (
	"github.com/home2mqtt/wscgo/config"
	"github.com/home2mqtt/wscgo/devices"
	"github.com/home2mqtt/wscgo/protocol"
)

type inputConfigPartParser struct{}

func (*inputConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.InputConfig{}
	section.FillData(s)
	c := protocol.CreateDInputConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(config.SLDevice, func(context config.RuntimeContext) error {
		device, err := devices.CreateInput(s)
		if err != nil {
			return err
		}
		context.AddDevice(device)
		context.AddProtocol(protocol.IntegrateInput(device, c))
		return nil
	})
	return nil
}

type outputConfigPartParser struct{}

func (*outputConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.OutputConfig{}
	section.FillData(s)
	c := protocol.CreateSwitchConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(config.SLDevice, func(context config.RuntimeContext) error {
		device, err := devices.CreateOutput(s)
		if err != nil {
			return err
		}
		context.AddDevice(device)
		context.AddProtocol(protocol.IntegrateSwitch(device, c))
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("dinput", &inputConfigPartParser{})
	config.RegisterConfigurationPartParser("switch", &outputConfigPartParser{})
}
