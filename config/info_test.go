package config

import (
	"bufio"
	"os"
	"testing"
)

func TestCpuInfo(t *testing.T) {
	cpuinfofile, _ := os.Open(cpuinfo)
	scanner := bufio.NewScanner(cpuinfofile)
	for scanner.Scan() {
		t.Log(scanner.Text())
	}
}

func TestModelInfo(t *testing.T) {
	model, serial, err := getModelInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(model)
	if model == "" {
		t.Fail()
	}
	t.Log(serial)
	if serial == "" {
		t.Fail()
	}
}

func TestDiscoveryInfo(t *testing.T) {
	version := "0.0.0-test"
	info := ComputeDeviceInfo(version)
	if info.Identifiers[0] == "" {
		t.Fail()
	}
	if info.Model == "" {
		t.Fail()
	}
	if info.SwVersion != version {
		t.Fail()
	}
	t.Log(info.Identifiers)
	t.Log(info.Model)
}
