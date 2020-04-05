package protocol

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grill-tamasi/wscgo/devices"
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
	BasicDiscoveryInfo
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

func (input *dinput) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	return &dinputDiscoveryConfig{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:       input.Name,
		StateTopic: input.StateTopic,
	}
}
