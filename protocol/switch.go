package protocol

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grill-tamasi/wscgo/devices"
	"periph.io/x/periph/conn/gpio"
)

type SwitchConfig struct {
	BasicDeviceConfig
	CommandTopic string `ini:"topic"`
}

type sw struct {
	devices.IOutput
	*SwitchConfig
}

//https://www.home-assistant.io/integrations/switch.mqtt/
type switchDiscoveryInfo struct {
	BasicDiscoveryInfo
	CommandTopic string `json:"command_topic,omitempty"`
	Name         string `json:"name,omitempty"`
}

func CreateSwitchConfig(id string) *SwitchConfig {
	return &SwitchConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
		},
	}
}

func IntegrateSwitch(output devices.IOutput, config *SwitchConfig) IDiscoverable {
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
			sw.Out(gpio.High)
		case "OFF":
			sw.Out(gpio.Low)
		default:
			value, err := strconv.Atoi(string(msg.Payload()))
			if err != nil {
				log.Println("WARNING: Switch ", sw.Name, " received unkown command: ", cmd)
			} else {
				if value > 0 {
					sw.Out(gpio.High)
				} else {
					sw.Out(gpio.Low)
				}
			}
		}
	})
}

func (sw *sw) GetComponent() string {
	return "switch"
}

func (sw *sw) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	return &switchDiscoveryInfo{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		CommandTopic: sw.CommandTopic,
		Name:         sw.Name,
	}
}
