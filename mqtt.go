package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttDevice interface {
	tick()
	init()
	configure(mqtt.Client)
}

func (serialconf *serialconf) configure(client mqtt.Client) {
	serialconf.valuecallback = func(arg int, value int) {
		switch arg {
		case 10:
			client.Publish(serialconf.topicroot+"/temperature", 0, false, fmt.Sprintf("%d", value))
		case 11:
			client.Publish(serialconf.topicroot+"/humidity", 0, false, fmt.Sprintf("%d", value))
		}
	}
	client.Subscribe(serialconf.topicroot+"/light", 0, func(client mqtt.Client, msg mqtt.Message) {
		value, err := strconv.Atoi(string(msg.Payload()))
		if err == nil {
			serialconf.setLight(value)
		}
	})
}

func (io *io) configure(client mqtt.Client) {
	client.Subscribe(io.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		value, err := strconv.Atoi(string(msg.Payload()))
		if err == nil {
			io.value = value
		}
	})
}

func (nand *nand) configure(client mqtt.Client) {

}

func (shutter *shutter) configure(client mqtt.Client) {
	// Setup single command topic
	rcvtopic := shutter.topic
	statustopic := rcvtopic + "/state"
	client.Subscribe(rcvtopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		value, err := strconv.Atoi(string(msg.Payload()))
		if err == nil {
			shutter.setCmd(value)
		}
	})
	log.Println("Topic (", rcvtopic, ") configured for shutter", shutter.id)

	// Callback config
	shutter.Callback = func(state int) {
		client.Publish(statustopic, 0, true, []byte(strconv.Itoa(state)))
	}

	// Setup group command topic
	groupRcvTopic := shutter.groupTopic
	if groupRcvTopic != "" {
		client.Subscribe(groupRcvTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
			value, err := strconv.Atoi(string(msg.Payload()))
			if err == nil {
				shutter.setCmd(value)
			}
		})
		log.Println("Group topic (", groupRcvTopic, ") configured for shutter", shutter.id)
	} else {
		log.Println("No group topic configured for shutter", shutter.id)
	}

	log.Println("Configured ", shutter, "{", rcvtopic)
}

func execute(options *mqtt.ClientOptions, devices []mqttDevice) {
	controlTicker := time.NewTicker(100 * time.Millisecond)

	client := mqtt.NewClient(options.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
		for _, p := range devices {
			p.configure(client)
			p.init()
		}
	}))
	log.Println("Connecting..")
	token := client.Connect()
	token.Wait()
	log.Println("Connected: ", token.Error())

	go func() {
		for range controlTicker.C {
			for _, p := range devices {
				p.tick()
			}
		}
	}()

	select {}
}
