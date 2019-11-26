package config

import (
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type DeviceInitializer func(devices.IoContext) protocol.IDiscoverable

type ConfigInitializer func(devices.IoContext)

type WscgoConfiguration struct {
	protocol.MqttConfig
	Configs []ConfigInitializer
	Devices []DeviceInitializer
}

func defaultConfiguration() *WscgoConfiguration {
	return &WscgoConfiguration{
		MqttConfig: protocol.MqttConfig{
			Host: "tcp://localhost:1883",
		},
	}
}
