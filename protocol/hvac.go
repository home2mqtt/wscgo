package protocol

import (
	"log"
	"strconv"
	"strings"

	"github.com/balazsgrill/wscgo/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// https://www.home-assistant.io/integrations/climate.mqtt/
type hvacDiscoveryInfo struct {
	BasicDiscoveryInfo
	Name string `json:"name,omitempty"`

	ActionTopic string `json:"action_topic,omitempty"`

	CurrentTemperatureTopic string `json:"current_temperature_topic,omitempty"`
	TemperatreCommandTopic  string `json:"temperature_command_topic,omitempty"`
	//TemperatureUnit         string  `json:"temperature_unit,omitempty"`
	TemperatureStateTopic string   `json:"temperature_state_topic,omitempty"`
	MaxTemp               float64  `json:"max_temp,omitempty"`
	MinTemp               float64  `json:"min_temp,omitempty"`
	TempStep              float64  `json:"temp_step,omitempty"`
	Modes                 []string `json:"modes,omitempty"`
	ModeCommandTopic      string   `json:"mode_command_topic,omitempty"`
}

// HVACConfig contains configuration parameters for a HVAC device
type HVACConfig struct {
	BasicDeviceConfig `ini:"Parent"`
	Topic             string `ini:"topic,omitempty"`
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

func (h *hvac) CurrentTemperatureTopic() string {
	return h.Topic + "/current"
}

func (h *hvac) TemperatureStateTopic() string {
	return h.Topic + "/state"
}

func (h *hvac) TemperatureCommandTopic() string {
	return h.Topic + "/set"
}

func (h *hvac) ModeCommandTopic() string {
	return h.Topic + "/mode"
}

func (h *hvac) onSetTemperature(client mqtt.Client, msg mqtt.Message) {
	cmd := string(msg.Payload())
	f, err := strconv.ParseFloat(cmd, 64)
	if err != nil {
		log.Printf("%s: Invalid target temperature value: '%s'", h.Name, cmd)
		return
	}
	h.SetTargetTemperature(f)
}

func (h *hvac) onSetMode(client mqtt.Client, msg mqtt.Message) {
	cmd := strings.ToUpper(string(msg.Payload()))
	r := h.TemperatureRange()
	switch cmd {
	case "OFF":
		h.SetTargetTemperature(r.Min)
	case "HEAT":
		h.SetTargetTemperature(r.Max)
	default:
		log.Printf("%s: Invalid command: '%s'", h.Name, cmd)
	}
}

func (h *hvac) Configure(client mqtt.Client) {
	ConfigureSensorListener(h.Temperature(), h.CurrentTemperatureTopic(), client)
	ConfigureSensorListener(h.TargetTemperature(), h.TemperatureStateTopic(), client)
	client.Subscribe(h.TemperatureCommandTopic(), 0, h.onSetTemperature)
	client.Subscribe(h.ModeCommandTopic(), 0, h.onSetMode)
}

func (h *hvac) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	r := h.TemperatureRange()
	return &hvacDiscoveryInfo{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:                    h.Name,
		CurrentTemperatureTopic: h.CurrentTemperatureTopic(),
		TemperatreCommandTopic:  h.TemperatureCommandTopic(),
		TemperatureStateTopic:   h.TemperatureStateTopic(),
		//TemperatureUnit:         "C",
		MaxTemp:          r.Max,
		MinTemp:          r.Min,
		TempStep:         r.Step,
		Modes:            []string{"off", "heat"},
		ModeCommandTopic: h.ModeCommandTopic(),
	}
}
