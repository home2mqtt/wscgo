package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

func execute(opts *mqtt.ClientOptions, conf *config.WscgoConfiguration) {
	client := mqtt.NewClient(opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
		for _, c := range conf.Configs {
			c
		}
	}))
	log.Println("Connecting..")
	token := client.Connect()
	token.Wait()
	log.Println("Connected: ", token.Error())

}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("This program expects a configuration file!")
	}
	conf := config.LoadConfig(args[0])

	opts := protocol.ConfigureClientOptions(conf.MqttConfig)

	execute(opts, conf)
}
