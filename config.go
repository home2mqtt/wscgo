package main

import (
	"log"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type deviceList struct {
	devices []mqttDevice
}

func getIntSafe(section *ini.Section, key string) int {
	if !section.HasKey(key) {
		log.Fatal("Required value ", key, " in ", section.Name())
	}
	k, err := section.GetKey(key)
	if err != nil {
		log.Fatal("Invalid value of ", key, " in ", section.Name(), ": ", err)
	}
	d, err := strconv.ParseInt(k.String(), 0, 32)
	if err != nil {
		log.Fatal("Invalid value of ", key, " in ", section.Name(), ": ", err)
	}
	return int(d)
}

func getStringSafe(section *ini.Section, key string) string {
	if !section.HasKey(key) {
		log.Fatal("Required value ", key, " in ", section.Name())
	}
	k, err := section.GetKey(key)
	if err != nil {
		log.Fatal("Invalid value of ", key, " in ", section.Name(), ": ", err)
	}
	return k.String()
}

func (devices *deviceList) processConfig(cat string, id string, section *ini.Section, ioContext ioContext) {
	switch cat {
	case "mcp23017":
		address := getIntSafe(section, "address")
		expansionBase := getIntSafe(section, "expansionBase")
		mcp23017Setup(expansionBase, address)
	case "shutter":
		topic := getStringSafe(section, "topic") + "/" + id
		devices.devices = append(devices.devices, &shutter{
			ioContext:     ioContext,
			topic:         topic,
			UpPin:         getIntSafe(section, "uppin"),
			DownPin:       getIntSafe(section, "downpin"),
			DirSwitchWait: getIntSafe(section, "dirswitchwait"),
			Range:         getIntSafe(section, "range"),
			PrevDir:       0,
		})
	case "nand":
		devices.devices = append(devices.devices, &nand{
			ioContext: ioContext,
			in1:       getIntSafe(section, "in1"),
			in2:       getIntSafe(section, "in2"),
			out:       getIntSafe(section, "out"),
		})
	case "io":
		devices.devices = append(devices.devices, &io{
			ioContext: ioContext,
			topic:     getStringSafe(section, "topic"),
			out:       getIntSafe(section, "out"),
		})
	case "serial":
		devices.devices = append(devices.devices, &serialconf{
			portname:  getStringSafe(section, "port"),
			baudrate:  uint(getIntSafe(section, "baud")),
			topicroot: getStringSafe(section, "topic"),
		})
	case ini.DEFAULT_SECTION:

	default:
		log.Fatalf("Unsupported configuration: %s", cat)
	}
}

func loadConfig(filename string) []mqttDevice {
	result := deviceList{
		devices: make([]mqttDevice, 0),
	}
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

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
		result.processConfig(category, id, s, &wiringPiIO{})
	}

	return result.devices
}
