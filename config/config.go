package config

import (
	"os"
	"sort"

	"github.com/balazsgrill/wscgo/devices"
	"github.com/balazsgrill/wscgo/protocol"
)

// DeviceInitializer is a function capable of initializing a device
type DeviceInitializer func(RuntimeContext) error

// RuntimeContext is holding devices and protocol conntections
type RuntimeContext interface {
	// AddDevice registers a device to this context
	AddDevice(devices.Device)
	// AddProtocol registers a protocol connection to this context
	AddProtocol(protocol.IDiscoverable)
}

// StartLevel is used to order in which the devices are initialized
type StartLevel int

// SLExtender is the default start level for IO extenders
const SLExtender StartLevel = 0

// SLDevice is the default start level for controlled devices
const SLDevice StartLevel = 10

// ConfigurationContext is holder initializers
type ConfigurationContext interface {
	AddDeviceInitializer(StartLevel, DeviceInitializer)
}

// WscgoConfiguration holds the configuration of a WSCGO instance
type WscgoConfiguration struct {
	protocol.MqttConfig
	devices []wscgoDeviceInitializerEntry

	// Node holds discovery metadata for the WSCGO instance
	Node protocol.DiscoverableNode
}

type wscgoDeviceInitializerEntry struct {
	startLevel        StartLevel
	deviceInitializer DeviceInitializer
}

// GetDeviceInitializers returns the list of initializers configured
func (config *WscgoConfiguration) GetDeviceInitializers() []DeviceInitializer {
	sort.SliceStable(config.devices, func(i int, j int) bool {
		return config.devices[i].startLevel < config.devices[j].startLevel
	})
	deviceInitializers := make([]DeviceInitializer, len(config.devices))
	for i, d := range config.devices {
		deviceInitializers[i] = d.deviceInitializer
	}
	return deviceInitializers
}

// AddDeviceInitializer registers a device initializer
func (config *WscgoConfiguration) AddDeviceInitializer(startLevel StartLevel, d DeviceInitializer) {
	config.devices = append(config.devices, wscgoDeviceInitializerEntry{
		startLevel:        startLevel,
		deviceInitializer: d,
	})
}

// ConfigurationSection denotes a section in the configuration
type ConfigurationSection interface {
	// FillData fills the values held by this instance to the given struct
	FillData(interface{}) error
	// GetID returns the locally unique identifier of this section
	GetID() string
}

// ConfigurationPartParser is an object capable of parsing a particular type of configuration section
type ConfigurationPartParser interface {
	// ParseConfiguration parses the given section expected to provide configuration the given context
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
