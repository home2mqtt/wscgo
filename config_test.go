package main

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := loadConfig("wscgo.ini")
	if config.MqttConfig.host != "tcp://localhost:1883" {
		t.Error("Mqtt broker host not configured properly")
	}
}
