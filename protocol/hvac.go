package protocol

import (
	"log"
	"strconv"
	"strings"

	"github.com/balazsgrill/hass"
	"github.com/balazsgrill/wscgo/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

func (h *hvac) GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig {
	r := h.TemperatureRange()
	return &hass.HVAC{
		BasicConfig: hass.BasicConfig{
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
