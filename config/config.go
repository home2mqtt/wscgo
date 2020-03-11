package config

import (
	"os"
	"plugin"
	"fmt"

	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type DeviceInitializer func() (protocol.IDiscoverable, error)

type ConfigInitializer func() error

type ConfigurationContext interface {
	AddConfigInitializer(ConfigInitializer)
	AddDeviceInitializer(DeviceInitializer)
}

type WscgoConfiguration struct {
	protocol.MqttConfig
	Node    protocol.DiscoverableNode
	Configs []ConfigInitializer
	Devices []DeviceInitializer
}

type WscgoPluginConfiguration struct {
	path string `ini:"path"`
}

func (config *WscgoConfiguration) AddConfigInitializer(c ConfigInitializer) {
	config.Configs = append(config.Configs, c)
}

func (config *WscgoConfiguration) AddDeviceInitializer(d DeviceInitializer) {
	config.Devices = append(config.Devices, d)
}

type ConfigurationSection interface {
	FillData(interface{}) error
	GetID() string
}

type ConfigurationPartParser interface {
	ParseConfiguration(ConfigurationSection, ConfigurationContext) error
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

func (pc *WscgoPluginConfiguration) Load() {
	_ err := plugin.Open(pc.path)
	if err != nil {
		fmt.Errorf("Could not load plugin %s: %v",pc.path, err)
	}
}
