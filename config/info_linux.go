// +build linux
package config

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const modelinfo string = "/proc/device-tree/model"
const cpuinfo string = "/proc/cpuinfo"

func getModelInfo() (string, string, error) {
	modelbytes, err := ioutil.ReadFile(modelinfo)
	var model string
	if err != nil {
		// Cannot read model file, assume generic linux
		model = "linux"
	}
	model = string(modelbytes)
	cpuinfofile, err := os.Open(cpuinfo)
	scanner := bufio.NewScanner(cpuinfofile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			keyvalue := strings.Split(line, ":")
			if strings.ToLower(strings.TrimSpace(keyvalue[0])) == "serial" {
				serialstring := strings.TrimSpace(keyvalue[1])
				return model, serialstring, nil
			}
		}
	}
	return "", "", errors.New("Invalid serial number")
}
