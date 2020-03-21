package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/config"
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"

	_ "gitlab.com/grill-tamasi/wscgo/integration"
)

type wscgoInstance struct {
	conf       *config.WscgoConfiguration
	client     mqtt.Client
	devices    []devices.Device
	protocols  []protocol.IDiscoverable
	deviceInfo *protocol.DeviceDiscoveryInfo
}

func (instance *wscgoInstance) AddDevice(device devices.Device) {
	instance.devices = append(instance.devices, device)
}

func (instance *wscgoInstance) AddProtocol(pro protocol.IDiscoverable) {
	log.Println("Initializing ", pro.GetComponent(), pro.GetObjectId())
	instance.protocols = append(instance.protocols, pro)
}

func (instance *wscgoInstance) intitializeDevices() {
	for _, c := range instance.conf.Configs {
		err := c()
		if err != nil {
			log.Println(err.Error())
		}
	}
	for _, d := range instance.conf.Devices {
		err := d(instance)
		if err != nil {
			log.Println(err.Error())
		}
	}
	for _, d := range instance.devices {
		err := d.Initialize()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (instance *wscgoInstance) eventOnConnected(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	for _, d := range instance.protocols {
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
	for _, p := range instance.protocols {
		protocol.PublisDiscoveryMessage(client, &instance.conf.Node, p, instance.deviceInfo)
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
				err := p.Tick()
				if err != nil {
					log.Println(err.Error())
				}
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
