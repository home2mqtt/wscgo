package protocol

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/wscgo/devices"
	"periph.io/x/conn/v3/gpio"
)

type LightConfig struct {
	BasicDeviceConfig
	CommandTopic string `ini:"topic,omitempty"`
}

type light struct {
	devices.IDimmer
	*LightConfig
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
		light.fireBrightnessEvent(client, gpio.DutyMax-1)
	case "OFF":
		light.Off()
		light.fireBrightnessEvent(client, 0)
	default:
		value, err := strconv.ParseInt(string(msg.Payload()), 10, 32)
		if err == nil {
			light.SetBrightness(gpio.Duty(value))
			light.fireBrightnessEvent(client, gpio.Duty(value))
		} else {
			log.Println("WARNING: Light ", light.Name, " received unkown command: ", cmd)
		}
	}
}

func (light *light) fireBrightnessEvent(client mqtt.Client, brightness gpio.Duty) {
	client.Publish(light.CommandTopic+"/brightness", 0, false, fmt.Sprintf("%d", brightness))
}

func (light *light) Configure(client mqtt.Client) {
	client.Subscribe(light.CommandTopic, 0, light.onMsgReceive)
}

func (light *light) GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig {
	return &hass.Light{
		BasicConfig: hass.BasicConfig{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:                   light.Name,
		CommandTopic:           light.CommandTopic,
		BrightnessCommandTopic: light.CommandTopic,
		BrightnessScale:        int32(gpio.DutyMax) - 1,
		BrightnessStateTopic:   light.CommandTopic + "/brightness",
		OnCommandType:          "brightness",
	}
}
