package protocol

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/wscgo/devices"
)

// SensorConfig is the protocol configuration of a sensor
type SensorConfig struct {
	BasicDeviceConfig
	UnitOfMeasurement string
	Topic             string
	Icon              string
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

// ConfigureSensorListener configures the sensor to publish a message on each measured sensor value
func ConfigureSensorListener(sensor devices.ISensor, topic string, client mqtt.Client) {
	sensor.AddListener(func(value float64) {
		log.Printf("Sending %f to %s\n", value, topic)
		t := client.Publish(topic, 0, false, fmt.Sprintf("%f", value))
		t.Wait()
		err := t.Error()
		if err != nil {
			log.Printf("Publish failed: %v\n", err)
		}
	})
}

func (sensor *sensor) Configure(client mqtt.Client) {
	ConfigureSensorListener(sensor, sensor.Topic, client)
}

func (sensor *sensor) GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig {
	return &hass.Sensor{
		BasicConfig: hass.BasicConfig{
			UniqueID: uniqueID,
			Device:   device,
		},
		Name:              sensor.Name,
		UnitOfMeasurement: sensor.Unit(),
		Topic:             sensor.Topic,
		Icon:              sensor.Icon,
	}
}
