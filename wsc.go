package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("This program expects a configuration file!")
	}
	config := loadConfig(args[0])

	opts := mqtt.NewClientOptions().AddBroker(config.MqttConfig.host).SetAutoReconnect(true)

	execute(opts, config.devices)
}
