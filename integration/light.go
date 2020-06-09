package integration

import (
	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/devices"
	"github.com/balazsgrill/wscgo/protocol"
)

type lightConfigurationParser struct {
}

func (*lightConfigurationParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.DimmerConfig{}
	if err := section.FillData(s); err != nil {
		return err
	}
	c := protocol.CreateLightConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(config.SLDevice, func(context config.RuntimeContext) error {
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
