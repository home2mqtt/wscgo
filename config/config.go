package config

import (
	"log"
	"os"
	"plugin"

	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/plugins"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type DeviceInitializer func(RuntimeContext) error

type ConfigInitializer func() error

type RuntimeContext interface {
	AddDevice(devices.Device)
	AddProtocol(protocol.IDiscoverable)
}

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
	Path string `ini:"path"`
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

func (pc *WscgoPluginConfiguration) Load() error {
	log.Printf("Loading %s\n", pc.Path)
	p, err := plugin.Open(pc.Path)
	if err != nil {
		return err
	}
	s, err := p.Lookup(plugins.AddonsGetterName)
	if err != nil {
		return err
	}
	f := s.(func() []plugins.Addon)
	addons := f()
	for _, a := range addons {
		err = loadAddon(a)
		if err != nil {
			log.Printf("Error loading addon from plugin: %v", err)
		}
	}
	return nil
}
