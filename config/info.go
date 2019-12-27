package config

import (
	"os"

	"gitlab.com/grill-tamasi/wscgo/protocol"
)

func ComputeDeviceInfo(version string) *protocol.DeviceDiscoveryInfo {
	model, serial, _ := getModelInfo()
	host, _ := os.Hostname()
	return &protocol.DeviceDiscoveryInfo{
		Identifiers:  []string{serial},
		Connections:  []string{},
		Manufacturer: "wscgo",
		Model:        model,
		Name:         "wscgo_" + host,
		SwVersion:    version,
	}
}
