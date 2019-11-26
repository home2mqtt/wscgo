package protocol

import (
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/devices"
)

type SwitchConfig struct {
	CommandTopic string
	Name         string
	ObjectId     string
}

type sw struct {
	devices.IOutput
	SwitchConfig
}

//https://www.home-assistant.io/integrations/switch.mqtt/
type switchDiscoveryInfo struct {
	CommandTopic string `json:"command_topic,omitempty"`
	Name         string `json:"name,omitempty"`
}

func IntegrateSwitch(output devices.IOutput, config SwitchConfig) IDiscoverable {
	return &sw{
		IOutput:      output,
		SwitchConfig: config,
	}
}

func (sw *sw) Configure(client mqtt.Client) {
	client.Subscribe(sw.CommandTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		cmd := strings.ToUpper(string(msg.Payload()))
		switch cmd {
		case "ON":
			sw.SetValue(true)
		case "OFF":
			sw.SetValue(false)
		default:
			log.Println("WARNING: Switch ", sw.Name, " received unkown command: ", cmd)
		}
	})
}

func (sw *sw) GetComponent() string {
	return "switch"
}
func (sw *sw) GetObjectId() string {
	return sw.ObjectId
}

func (sw *sw) GetDiscoveryInfo() interface{} {
	return &switchDiscoveryInfo{
		CommandTopic: sw.CommandTopic,
		Name:         sw.Name,
	}
}
