package integration

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/wscgo/config"
	"github.com/home2mqtt/wscgo/devices"
	"github.com/home2mqtt/wscgo/protocol"
)

type WscgoInstance struct {
	Version string

	client     mqtt.Client
	devices    []devices.Device
	protocols  []protocol.IDiscoverable
	deviceInfo *hass.Device
	nodeInfo   *protocol.DiscoverableNode
}

func (instance *WscgoInstance) AddDevice(device devices.Device) {
	instance.devices = append(instance.devices, device)
}

func (instance *WscgoInstance) AddProtocol(pro protocol.IDiscoverable) {
	log.Println("Initializing ", pro.GetObjectId())
	instance.protocols = append(instance.protocols, pro)
}

func (instance *WscgoInstance) Configure(conf *config.WscgoConfiguration) {
	instance.nodeInfo = &conf.Node
	instance.SetMqttClientOptions(&conf.MqttConfig)
	for _, d := range conf.GetDeviceInitializers() {
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

func (instance *WscgoInstance) eventOnConnected(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	for _, d := range instance.protocols {
		log.Println("Configuring ", d.GetObjectId())
		d.Configure(client)
	}
	instance.publishDiscoveryMessages(client)
	/* publish anything to topic 'discover' to trigger discovery messages */
	client.Subscribe("discover", 0, func(client mqtt.Client, _ mqtt.Message) {
		instance.publishDiscoveryMessages(client)
	})
}

func (instance *WscgoInstance) publishDiscoveryMessages(client mqtt.Client) {
	for _, p := range instance.protocols {
		protocol.PublisDiscoveryMessage(client, instance.nodeInfo, p, instance.deviceInfo)
	}
}

func (instance *WscgoInstance) eventOnDisconnected(client mqtt.Client, err error) {
	log.Println("MQTT Connection lost " + err.Error())
}

func (instance *WscgoInstance) SetMqttClientOptions(conf *protocol.MqttConfig) {
	if instance.client != nil {
		instance.client.Disconnect(100)
	}
	opts := protocol.ConfigureClientOptions(conf)
	opts = opts.SetOnConnectHandler(instance.eventOnConnected).SetConnectionLostHandler(instance.eventOnDisconnected)
	instance.deviceInfo = config.ComputeDeviceInfo(instance.Version)
	instance.client = mqtt.NewClient(opts)
}

func (instance *WscgoInstance) Loop() {
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

func (instance *WscgoInstance) Start() {
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
