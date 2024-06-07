package protocol

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/wscgo/devices"
)

// CoverConfig is the protocol configuration of a cover (a.k.a. shutter)
type CoverConfig struct {
	BasicDeviceConfig `ini:"Parent"`
	CommandTopic      string `ini:"topic,omitempty"`
	GroupTopic        string `ini:"opt_groupTopic"`
	PositionTopic     string `ini:"position_topic,omitempty"`
}

type cover struct {
	devices.IShutter
	*CoverConfig
}

// CreateCoverConfig creates a CoverConfig structure with default values
func CreateCoverConfig(id string) *CoverConfig {
	return &CoverConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

// IntegrateCover initializes protocol on top of the given device
func IntegrateCover(shutter devices.IShutter, config *CoverConfig) IDiscoverable {
	return &cover{
		IShutter:    shutter,
		CoverConfig: config,
	}
}

func (cover *cover) GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig {
	return &hass.Cover{
		BasicConfig: hass.BasicConfig{
			UniqueID: uniqueID,
			Device:   device,
		},
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
	case "OPENORSTOP":
		cover.OpenOrStop()
	case "CLOSEORSTOP":
		cover.CloseOrStop()
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
		client.Publish(cover.PositionTopic, 0, true, strconv.Itoa(value))
	})
}
