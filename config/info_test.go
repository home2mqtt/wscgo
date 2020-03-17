package config

import (
	"encoding/json"
	"testing"
)

func TestModelInfo(t *testing.T) {
	model, serial, err := getModelInfo()
	if err != nil {
		t.Fatal(err)
	}
	if model == "" {
		t.Error("No model identification")
	}
	if serial == "" {
		t.Error("No serial number")
	}
}

func TestDiscoveryInfo(t *testing.T) {
	version := "0.0.0-test"
	info := ComputeDeviceInfo(version)
	if info.Identifiers[0] == "" {
		t.Error("No identifier")
	}
	if info.Model == "" {
		t.Error("No model name")
	}
	if info.SwVersion != version {
		t.Error(info.SwVersion)
	}
	_, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
	}
}
