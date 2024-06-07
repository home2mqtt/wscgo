package protocol

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/wscgo/devices"
)

type DInputConfig struct {
	BasicDeviceConfig
	StateTopic string `ini:"topic,omitempty"`
}

type dinput struct {
	devices.IInput
	*DInputConfig
}

func CreateDInputConfig(id string) *DInputConfig {
	return &DInputConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

func IntegrateInput(input devices.IInput, conf *DInputConfig) IDiscoverable {
	return &dinput{
		IInput:       input,
		DInputConfig: conf,
	}
}

func (input *dinput) Configure(client mqtt.Client) {
	input.AddListener(func(state bool) {
		value := "OFF"
		if state {
			value = "ON"
		}
		client.Publish(input.StateTopic, 0, true, value)
	})
}

func (input *dinput) GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig {
	return &hass.DInput{
		BasicConfig: hass.BasicConfig{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:       input.Name,
		StateTopic: input.StateTopic,
	}
}
