package config

import (
	"testing"
)

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
