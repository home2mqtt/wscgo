package integration

import (
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type shutterConfigPartParser struct{}

func (*shutterConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.ShutterConfig{}
	section.FillData(s)
	c := protocol.CreateCoverConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(func() (protocol.IDiscoverable, error) {
		shutter, err := devices.CreateShutter(s)
		if err != nil {
			return nil, err
		}
		return protocol.IntegrateCover(shutter, c), nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("shutter", &shutterConfigPartParser{})
}
