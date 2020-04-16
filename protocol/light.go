package protocol

import (
	"log"
	"strconv"
	"strings"

	"github.com/balazsgrill/wscgo/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LightConfig struct {
	BasicDeviceConfig
	CommandTopic string `ini:"topic,omitempty"`
}

type light struct {
	devices.IDimmer
	*LightConfig
}

type lightDiscoveryInfo struct {
	BasicDiscoveryInfo
	CommandTopic           string `json:"command_topic,omitempty"`
	Name                   string `json:"name,omitempty"`
	BrightnessCommandTopic string `json:"brightness_command_topic,omitempty"`
	BrightnessScale        int    `json:"brightness_scale"`
	BrightnessStateTopic   string `json:"brightness_state_topic,omitempty"`
	OnCommandType          string `json:"on_command_type,omitempty"`
	StateTopic             string `json:"state_topic,omitempty"`
}

func CreateLightConfig(id string) *LightConfig {
	return &LightConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

func IntegrateLight(dimmer devices.IDimmer, config *LightConfig) IDiscoverable {
	return &light{
		IDimmer:     dimmer,
		LightConfig: config,
	}
}

func (light *light) onMsgReceive(client mqtt.Client, msg mqtt.Message) {
	cmd := strings.ToUpper(string(msg.Payload()))
	switch cmd {
	case "ON":
		light.On()
		light.fireBrightnessEvent(client, light.BrightnessResolution()-1)
	case "OFF":
		light.Off()
		light.fireBrightnessEvent(client, 0)
	default:
		value, err := strconv.Atoi(string(msg.Payload()))
		if err == nil {
			light.SetBrightness(value)
			light.fireBrightnessEvent(client, value)
		} else {
			log.Println("WARNING: Light ", light.Name, " received unkown command: ", cmd)
		}
	}
}

func (light *light) fireBrightnessEvent(client mqtt.Client, brightness int) {
	client.Publish(light.CommandTopic+"/brightness", 0, false, strconv.Itoa(brightness))
}

func (light *light) GetComponent() string {
	return "light"
}

func (light *light) Configure(client mqtt.Client) {
	client.Subscribe(light.CommandTopic, 0, light.onMsgReceive)
}

func (light *light) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	return &lightDiscoveryInfo{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:                   light.Name,
		CommandTopic:           light.CommandTopic,
		BrightnessCommandTopic: light.CommandTopic,
		BrightnessScale:        light.BrightnessResolution() - 1,
		BrightnessStateTopic:   light.CommandTopic + "/brightness",
		OnCommandType:          "brightness",
	}
}
