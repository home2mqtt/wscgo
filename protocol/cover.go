package protocol

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/devices"
)

type CoverConfig struct {
	BasicDeviceConfig `ini:"Parent"`
	CommandTopic      string `ini:"topic,omitempty"`
	GroupTopic        string `ini:"opt_groupTopic"`
	PositionTopic     string `ini:"position_topic,omitempty"`
}

//https://www.home-assistant.io/integrations/cover.mqtt/
type coverDiscoveryInfo struct {
	CommandTopic   string `json:"command_topic,omitempty"`
	Name           string `json:"name,omitempty"`
	PositionTopic  string `json:"position_topic,omitempty"`
	PositionOpen   int    `json:"position_open"`
	PositionClosed int    `json:"position_closed"`
}

type cover struct {
	devices.IShutter
	*CoverConfig
}

func CreateCoverConfig(id string) *CoverConfig {
	return &CoverConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

func IntegrateCover(shutter devices.IShutter, config *CoverConfig) IDiscoverable {
	return &cover{
		IShutter:    shutter,
		CoverConfig: config,
	}
}

func (cover *cover) GetDiscoveryInfo() interface{} {
	return &coverDiscoveryInfo{
		CommandTopic:   cover.CommandTopic,
		Name:           cover.Name,
		PositionTopic:  cover.PositionTopic,
		PositionOpen:   cover.GetRange(),
		PositionClosed: 0,
	}
}

func (cover *cover) onMsgReceive(client mqtt.Client, msg mqtt.Message) {
	cmd := strings.ToUpper(string(msg.Payload()))
	switch cmd {
	case "OPEN":
		cover.Open()
	case "CLOSE":
		cover.Close()
	case "STOP":
		cover.Stop()
	default:
		value, err := strconv.Atoi(string(msg.Payload()))
		if err == nil {
			cover.MoveBy(value)
		} else {
			log.Println("WARNING: Cover ", cover.Name, " received unkown command: ", cmd)
		}
	}
}

func (cover *cover) Configure(client mqtt.Client) {
	client.Subscribe(cover.CommandTopic, 0, cover.onMsgReceive)
	if cover.GroupTopic != "" {
		client.Subscribe(cover.GroupTopic, 0, cover.onMsgReceive)
	}
	cover.AddListener(func(value int) {
		client.Publish(cover.PositionTopic, 0, false, strconv.Itoa(value))
	})
}

func (cover *cover) GetComponent() string {
	return "cover"
}
