package main

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/protocol"
	"gitlab.com/grill-tamasi/wscgo/tests"
)

func main() {
	io := tests.CreateTestIo(2)
	shutterConfig := &devices.ShutterConfig{
		UpPin:         0,
		DownPin:       1,
		Range:         50,
		DirSwitchWait: 10,
	}
	coverConfig := &protocol.CoverConfig{
		CommandTopic:  "test/wscgo/shutter2/cmd",
		Name:          "Test Shutter",
		PositionTopic: "test/wscgo/shutter2/pos",
		ObjectId:      "a1234",
	}
	mqttc := &protocol.MqttConfig{
		Host: "tcp://192.168.0.1:1883",
	}

	shutter := devices.CreateShutter(io, shutterConfig)
	cover := protocol.IntegrateCover(shutter, coverConfig)

	clientOpts := protocol.ConfigureClientOptions(mqttc)
	clientOpts.SetOnConnectHandler(func(client mqtt.Client) {
		cover.Configure(client)
		protocol.PublisDiscoveryMessage(client, &protocol.DiscoverableNode{
			DiscoveryPrefix: "protocol",
			NodeID:          "DiscoveryTest",
		}, cover)
	})

	client := mqtt.NewClient(clientOpts)
	token := client.Connect()
	token.Wait()

	controlTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range controlTicker.C {
			shutter.Tick()
		}
	}()

	select {}
}
