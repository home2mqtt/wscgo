package protocol

import (
	"log"
	"strconv"

	"github.com/balazsgrill/wscgo/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// https://www.home-assistant.io/integrations/climate.mqtt/
type hvacDiscoveryInfo struct {
	BasicDiscoveryInfo
	Name string `json:"name,omitempty"`

	ActionTopic string `json:"action_topic,omitempty"`

	CurrentTemperatureTopic string  `json:"current_temperature_topic,omitempty"`
	TemperatreCommandTopic  string  `json:"temperature_command_topic,omitempty"`
	TemperatureUnit         string  `json:"temperature_unit,omitempty"`
	TemperatureStateTopic   string  `json:"temperature_state_topic,omitempty"`
	MaxTemp                 float64 `json:"max_temp,omitempty"`
	MinTemp                 float64 `json:"min_temp,omitempty"`
	TempStep                float64 `json:"temp_step,omitempty"`
}

// HVACConfig contains configuration parameters for a HVAC device
type HVACConfig struct {
	BasicDeviceConfig       `ini:"Parent"`
	CurrentTemperatureTopic string `json:"current_temperature_topic,omitempty"`
	TemperatreCommandTopic  string `json:"temperature_command_topic,omitempty"`
	TemperatureStateTopic   string `json:"temperature_state_topic,omitempty"`
}

type hvac struct {
	devices.IThermostat
	*HVACConfig
}

func CreateHVACConfig(id string) *HVACConfig {
	return &HVACConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

func IntegrateHVAC(device devices.IThermostat, config *HVACConfig) IDiscoverable {
	return &hvac{
		IThermostat: device,
		HVACConfig:  config,
	}
}

func (h *hvac) GetComponent() string {
	return "climate"
}

func (h *hvac) onMsgReceive(client mqtt.Client, msg mqtt.Message) {
	cmd := string(msg.Payload())
	f, err := strconv.ParseFloat(cmd, 64)
	if err != nil {
		log.Printf("%s: Invalid target temperature value: '%s'", h.Name, cmd)
		return
	}
	h.SetTargetTemperature(f)
}

func (h *hvac) Configure(client mqtt.Client) {
	ConfigureSensorListener(h.Temperature(), h.CurrentTemperatureTopic, client)
	ConfigureSensorListener(h.TargetTemperature(), h.TemperatureStateTopic, client)
	client.Subscribe(h.TemperatreCommandTopic, 0, h.onMsgReceive)
}

func (h *hvac) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	r := h.TemperatureRange()
	return &hvacDiscoveryInfo{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:                    h.Name,
		CurrentTemperatureTopic: h.CurrentTemperatureTopic,
		TemperatreCommandTopic:  h.TemperatreCommandTopic,
		TemperatureStateTopic:   h.TemperatureStateTopic,
		TemperatureUnit:         "C",
		MaxTemp:                 r.Max,
		MinTemp:                 r.Min,
		TempStep:                r.Step,
	}
}
