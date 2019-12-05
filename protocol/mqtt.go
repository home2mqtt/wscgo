package protocol

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/devices"
)

type MQTTProtocol interface {
	devices.Device
	Configure(mqtt.Client)
}

type MqttConfig struct {
	// Host, e.g. tcp://localhost:1883
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	ClientID string `ini:"clientid"`
}

func ConfigureClientOptions(config *MqttConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(config.Host).SetAutoReconnect(true)
	if config.User != "" {
		opts = opts.SetUsername(config.User)
	}
	if config.Password != "" {
		opts = opts.SetPassword(config.Password)
	}

	return opts
}
