package integration

import (
	"log"
	"time"

	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/devices"
	"github.com/balazsgrill/wscgo/protocol"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type WscgoInstance struct {
	client     mqtt.Client
	devices    []devices.Device
	protocols  []protocol.IDiscoverable
	deviceInfo *protocol.DeviceDiscoveryInfo
	nodeInfo   *protocol.DiscoverableNode
}

func (instance *WscgoInstance) AddDevice(device devices.Device) {
	instance.devices = append(instance.devices, device)
}

func (instance *WscgoInstance) AddProtocol(pro protocol.IDiscoverable) {
	log.Println("Initializing ", pro.GetComponent(), pro.GetObjectId())
	instance.protocols = append(instance.protocols, pro)
}

func (instance *WscgoInstance) Configure(conf *config.WscgoConfiguration) {
	instance.SetMqttClientOptions(protocol.ConfigureClientOptions(&conf.MqttConfig))
	instance.nodeInfo = &conf.Node
	for _, c := range conf.Configs {
		err := c()
		if err != nil {
			log.Println(err.Error())
		}
	}
	for _, d := range conf.Devices {
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
		log.Println("Configuring ", d.GetComponent(), d.GetObjectId())
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

func (instance *WscgoInstance) eventOnDisconnected(client mqtt.Client) {
	log.Println("MQTT Connection lost")
}

func (instance *WscgoInstance) SetMqttClientOptions(opts *mqtt.ClientOptions) {
	if instance.client != nil {
		instance.client.Disconnect(100)
	}
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
