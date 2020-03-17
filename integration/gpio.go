package integration

import (
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type inputConfigPartParser struct{}

func (*inputConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.InputConfig{}
	section.FillData(s)
	c := protocol.CreateDInputConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(func() (protocol.IDiscoverable, error) {
		device, err := devices.CreateInput(s)
		if err != nil {
			return nil, err
		}
		return protocol.IntegrateInput(device, c), nil
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
	context.AddDeviceInitializer(func() (protocol.IDiscoverable, error) {
		device, err := devices.CreateOutput(s)
		if err != nil {
			return nil, err
		}
		return protocol.IntegrateSwitch(device, c), nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("dinput", &inputConfigPartParser{})
	config.RegisterConfigurationPartParser("switch", &outputConfigPartParser{})
}
