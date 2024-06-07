//go:build linux
// +build linux

package integration

import (
	"github.com/home2mqtt/wscgo/config"
	"github.com/home2mqtt/wscgo/devices/cometblue"
	"github.com/home2mqtt/wscgo/protocol"
)

type cometblueConfigurationParser struct {
}

func (*cometblueConfigurationParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	s := &cometblue.Config{}
	if err := section.FillData(s); err != nil {
		return err
	}
	c := protocol.CreateHVACConfig(section.GetID())
	section.FillData(&c.BasicDeviceConfig)
	section.FillData(c)
	context.AddDeviceInitializer(config.SLDevice, func(context config.RuntimeContext) error {
		device := cometblue.Create(s)
		context.AddDevice(device)
		context.AddProtocol(protocol.IntegrateHVAC(device, c))
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("cometblue", &cometblueConfigurationParser{})
}
