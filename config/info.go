package config

import (
	"os"

	"github.com/balazsgrill/wscgo/protocol"
)

// ComputeDeviceInfo extracts discovery metadata from the host system
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
