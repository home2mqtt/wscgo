package protocol

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/devices"
)

type DInputConfig struct {
	BasicDeviceConfig
	StateTopic string `ini:"topic,omitempty"`
}

type dinput struct {
	devices.IInput
	*DInputConfig
}

type dinputDiscoveryConfig struct {
	Name       string `json:"name"`
	StateTopic string `json:"state_topic"`
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

func (input *dinput) GetComponent() string {
	return "binary_sensor"
}

func (input *dinput) GetDiscoveryInfo() interface{} {
	return &dinputDiscoveryConfig{
		Name:       input.Name,
		StateTopic: input.StateTopic,
	}
}
