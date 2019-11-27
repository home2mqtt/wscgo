package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/protocol"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

const version string = "0.2-beta1"

type wscgoInstance struct {
	conf    *config.WscgoConfiguration
	client  mqtt.Client
	devices []protocol.IDiscoverable
}

func (instance *wscgoInstance) intitializeDevices() {
	for _, c := range instance.conf.Configs {
		c(instance.conf.IoContext)
	}
	for _, d := range instance.conf.Devices {
		instance.devices = append(instance.devices, d(instance.conf.IoContext))
	}
	for _, d := range instance.devices {
		d.Initialize()
	}
}

func (instance *wscgoInstance) eventOnConnected(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	for _, d := range instance.devices {
		d.Configure(client)
		protocol.PublisDiscoveryMessage(client, &instance.conf.Node, d)
	}
}

func (instance *wscgoInstance) eventOnDisconnected(client mqtt.Client) {
	log.Println("MQTT Connection lost")
}

func CreateWscgo(conf *config.WscgoConfiguration) *wscgoInstance {
	instance := &wscgoInstance{
		conf: conf,
	}
	opts := protocol.ConfigureClientOptions(&conf.MqttConfig)
	opts = opts.SetOnConnectHandler(instance.eventOnConnected)
	instance.client = mqtt.NewClient(opts)
	return instance
}

func (instance *wscgoInstance) loop() {
	controlTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range controlTicker.C {
			for _, p := range instance.devices {
				p.Tick()
			}
		}
	}()

	select {}
}

func (instance *wscgoInstance) start() {
	log.Println("Attempting to connect to MQTT broker")
	for !instance.client.IsConnected() {
		token := instance.client.Connect()
		token.Wait()
		err := token.Error()
		if err != nil {
			log.Println("Connection failed: ", token.Error())
		}
	}
}

func main() {
	log.Println("wscgo version ", version)
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("usage: wscgo config.ini")
	}
	conf := config.LoadConfig(args[0])
	conf.IoContext = &wiringpi.WiringPiIO{}

	instance := CreateWscgo(conf)
	instance.intitializeDevices()
	instance.start()
	instance.loop()
}
