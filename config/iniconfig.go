package config

import (
	"log"
	"strings"

	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
	"gopkg.in/ini.v1"
)

func (conf *WscgoConfiguration) processConfig(category string, id string, section *ini.Section) {
	switch category {
	case "mqtt":
		section.MapTo(conf.MqttConfig)
	case "mcp23017":
		c := &wiringpi.Mcp23017Config{}
		section.MapTo(c)
		conf.Configs = append(conf.Configs, func(devices.IoContext) {
			wiringpi.Mcp23017Setup(c)
		})
	case "shutter":
		s := &devices.ShutterConfig{}
		section.MapTo(s)
		c := protocol.CreateCoverConfig(id)
		section.MapTo(c)
		conf.Devices = append(conf.Devices, func(io devices.IoContext) protocol.IDiscoverable {
			shutter := devices.CreateShutter(io, s)
			return protocol.IntegrateCover(shutter, c)
		})
	case "switch":
		s := &devices.OutputConfig{}
		section.MapTo(s)
		c := protocol.CreateSwitchConfig(id)
		section.MapTo(c)
		conf.Devices = append(conf.Devices, func(io devices.IoContext) protocol.IDiscoverable {
			device := devices.CreateOutput(io, s)
			return protocol.IntegrateSwitch(device, c)
		})
	}
}

func LoadConfig(filename string) *WscgoConfiguration {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	conf := defaultConfiguration()

	for _, s := range cfg.Sections() {
		name := s.Name()
		cat := strings.Split(name, ":")
		l := len(cat)
		var category string
		var id string
		category = cat[0]
		id = ""
		if l > 1 {
			id = cat[1]
		}
		conf.processConfig(category, id, s)
	}

	return conf
}
