package config

import (
	"log"
	"strings"

	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/homeassistant"
	"gitlab.com/grill-tamasi/wscgo/protocol"
	"gopkg.in/ini.v1"
)

type DeviceInitializer func(devices.IoContext) homeassistant.IDiscoverable

type wscgoConfiguration struct {
	protocol.MqttConfig
	devices []DeviceInitializer
}

func (conf *wscgoConfiguration) processConfig(category string, id string, section *ini.Section) {
	switch category {
	case "shutter":
		s := &devices.ShutterConfig{}
		section.MapTo(s)
		c := &homeassistant.CoverConfig{}
		section.MapTo(c)
		conf.devices = append(conf.devices, func(io devices.IoContext) homeassistant.IDiscoverable {
			shutter := devices.CreateShutter(io, s)
			return homeassistant.IntegrateCover(shutter, c)
		})
	}
}

func LoadConfig(filename string) {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	conf := &wscgoConfiguration{}

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
}
