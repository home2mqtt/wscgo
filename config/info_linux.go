// +build linux

package config

import (
	"io/ioutil"

	"periph.io/x/host/v3/distro"
)

const modelinfo string = "/proc/device-tree/model"
const cpuinfo string = "/proc/cpuinfo"

func getModelInfo() (string, string, error) {
	modelbytes, err := ioutil.ReadFile(modelinfo)
	var model string
	if err != nil {
		// Cannot read model file, assume generic linux
		model = distro.OSRelease()["PRETTY_NAME"]
	} else {
		model = string(modelbytes)
	}

	var serial string
	var ok bool
	if serial, ok = distro.CPUInfo()["serial"]; !ok {
		serial = "Unkown"
	}
	return model, serial, nil
}
