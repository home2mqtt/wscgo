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
		t.Fail()
	}
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
		t.Error(info.SwVersion)
	}
	_, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
	}
}
