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
	Version string

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

func (instance *WscgoInstance) SetMqttClientOptions(conf *protocol.MqttConfig) {
	if instance.client != nil {
		instance.client.Disconnect(100)
	}
	opts := protocol.ConfigureClientOptions(conf)
	opts = opts.SetOnConnectHandler(instance.eventOnConnected)
	instance.client = mqtt.NewClient(opts)
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
