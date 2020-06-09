package protocol

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTProtocol interface {
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
	log.Println("Configured MQTT broker: ", config.Host)
	opts := mqtt.NewClientOptions().AddBroker(config.Host).SetAutoReconnect(true)
	if config.User != "" {
		opts = opts.SetUsername(config.User)
	}
	if config.Password != "" {
		opts = opts.SetPassword(config.Password)
	}

	return opts
}
