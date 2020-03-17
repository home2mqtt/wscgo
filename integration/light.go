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
	context.AddDeviceInitializer(func() (protocol.IDiscoverable, error) {
		device, err := devices.CreateDimmer(s)
		if err != nil {
			return nil, err
		}
		return protocol.IntegrateLight(device, c), nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("light", &lightConfigurationParser{})
}
