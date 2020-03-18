package integration

import (
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type lightConfigurationParser struct {
}

func (*lightConfigurationParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.DimmerConfig{
		Resolution: 1024,
	}
	if err := section.FillData(s); err != nil {
		return err
	}
	c := protocol.CreateLightConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(func(context config.RuntimeContext) error {
		device, err := devices.CreateDimmer(s)
		if err != nil {
			return err
		}
		context.AddDevice(device)
		context.AddProtocol(protocol.IntegrateLight(device, c))
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("light", &lightConfigurationParser{})
}
