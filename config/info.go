package config

import (
	"os"

	"github.com/home2mqtt/hass"
)

// ComputeDeviceInfo extracts discovery metadata from the host system
func ComputeDeviceInfo(version string) *hass.Device {
	model, serial, _ := getModelInfo()
	host, _ := os.Hostname()
	return &hass.Device{
		Identifiers:  []string{serial},
		Connections:  []string{},
		Manufacturer: "wscgo",
		Model:        model,
		Name:         "wscgo_" + host,
		SwVersion:    version,
	}
}
