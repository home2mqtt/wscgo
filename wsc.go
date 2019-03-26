package main

import "github.com/eclipse/paho.mqtt.golang"

import "os"
import "log"

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("This program expects a configuration file!")
	}
	devices := loadConfig(args[0])

	opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.0.1:1883").SetClientID("test0").SetAutoReconnect(true)

	execute(opts, devices)
}
