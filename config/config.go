package config

import (
	"os"

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
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "wscgo"
	}
	return &WscgoConfiguration{
		Node: protocol.DiscoverableNode{
			DiscoveryPrefix: "homeassistant",
			NodeID:          hostname,
		},
		MqttConfig: protocol.MqttConfig{
			Host: "tcp://localhost:1883",
		},
	}
}
