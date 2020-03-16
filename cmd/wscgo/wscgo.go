package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/protocol"

	_ "gitlab.com/grill-tamasi/wscgo/integration"
)

type wscgoInstance struct {
	conf       *config.WscgoConfiguration
	client     mqtt.Client
	devices    []protocol.IDiscoverable
	deviceInfo *protocol.DeviceDiscoveryInfo
}

func (instance *wscgoInstance) intitializeDevices() {
	for _, c := range instance.conf.Configs {
		c()
	}
	for _, d := range instance.conf.Devices {
		dev, err := d()
		if err != nil {
			log.Println(err.Error())
		} else {
			instance.devices = append(instance.devices, dev)
		}
	}
	for _, d := range instance.devices {
		log.Println("Initializing ", d.GetComponent(), d.GetObjectId())
		d.Initialize()
	}
}

func (instance *wscgoInstance) eventOnConnected(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	for _, d := range instance.devices {
		log.Println("Configuring ", d.GetComponent(), d.GetObjectId())
		d.Configure(client)
	}
	instance.publishDiscoveryMessages(client)
	/* publish anything to topic 'discover' to trigger discovery messages */
	client.Subscribe("discover", 0, func(client mqtt.Client, _ mqtt.Message) {
		instance.publishDiscoveryMessages(client)
	})
}

func (instance *wscgoInstance) publishDiscoveryMessages(client mqtt.Client) {
	for _, d := range instance.devices {
		protocol.PublisDiscoveryMessage(client, &instance.conf.Node, d, instance.deviceInfo)
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
	instance.deviceInfo = config.ComputeDeviceInfo(Version)
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
	log.Println("wscgo started")
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
	log.Println("wscgo version ", Version)
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("usage: wscgo config.ini")
	}
	conf := config.LoadConfig(args[0])

	instance := CreateWscgo(conf)
	instance.intitializeDevices()
	instance.start()
	instance.loop()
}
