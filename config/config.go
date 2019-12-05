package config

import (
	"gitlab.com/grill-tamasi/wscgo/protocol"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

type DeviceInitializer func(wiringpi.IoContext) protocol.IDiscoverable

type ConfigInitializer func(wiringpi.IoContext)

type WscgoConfiguration struct {
	protocol.MqttConfig
	wiringpi.IoContext
	Node    protocol.DiscoverableNode
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
