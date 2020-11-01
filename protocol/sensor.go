package protocol

import (
	"fmt"

	"github.com/balazsgrill/wscgo/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SensorConfig is the protocol configuration of a sensor
type SensorConfig struct {
	BasicDeviceConfig
	UnitOfMeasurement string
	Topic             string
	Icon              string
}

//https://www.home-assistant.io/integrations/sensor.mqtt/
type sensorDiscoveryInfo struct {
	BasicDiscoveryInfo
	Name              string `json:"name,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	Topic             string `json:"state_topic"`
	Icon              string `json:"icon"`
}

type sensor struct {
	devices.ISensor
	*SensorConfig
}

// CreateSensorConfig provides the defaule configuration values for a sensor device
func CreateSensorConfig(id string) *SensorConfig {
	return &SensorConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			ObjectId: id,
			Name:     id,
		},
	}
}

// IntegrateSensor creates integration of the given sensor
func IntegrateSensor(sensordev devices.ISensor, config *SensorConfig) IDiscoverable {
	return &sensor{
		ISensor:      sensordev,
		SensorConfig: config,
	}
}

func (sensor *sensor) GetComponent() string {
	return "sensor"
}

// ConfigureSensorListener configures the sensor to publish a message on each measured sensor value
func ConfigureSensorListener(sensor devices.ISensor, topic string, client mqtt.Client) {
	sensor.AddListener(func(value float64) {
		client.Publish(topic, 0, false, fmt.Sprintf("%g", value))
	})
}

func (sensor *sensor) Configure(client mqtt.Client) {
	ConfigureSensorListener(sensor, sensor.Topic, client)
}

func (sensor *sensor) GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{} {
	return &sensorDiscoveryInfo{
		BasicDiscoveryInfo: BasicDiscoveryInfo{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:              sensor.Name,
		UnitOfMeasurement: sensor.Unit(),
		Topic:             sensor.Topic,
		Icon:              sensor.Icon,
	}
}
