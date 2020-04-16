package integration

import (
	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/devices"
	"github.com/balazsgrill/wscgo/protocol"
)

type shutterConfigPartParser struct{}

func (*shutterConfigPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &devices.ShutterConfig{}
	section.FillData(s)
	c := protocol.CreateCoverConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(func(context config.RuntimeContext) error {
		shutter, err := devices.CreateShutter(s)
		if err != nil {
			return err
		}
		context.AddDevice(shutter)
		context.AddProtocol(protocol.IntegrateCover(shutter, c))
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("shutter", &shutterConfigPartParser{})
}
