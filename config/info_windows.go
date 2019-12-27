// +build windows
package config

import (
	"errors"
	"os/exec"
	"strings"
)

func getModelInfo() (string, string, error) {
	cmd := exec.Command("wmic", "bios", "get", "serialnumber")
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	lines := strings.Split(string(output), "\r\n")
	for _, aline := range lines {
		line := strings.TrimSpace(aline)
		if line != "SerialNumber" && line != "" {
			return "windows", line, nil
		}
	}
	return "", "", errors.New("Invalid serial number")
}
