package config

import (
	"testing"
	"encoding/json"
)

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
	data, _ := json.Marshal(info)
	t.Log(string(data))
}