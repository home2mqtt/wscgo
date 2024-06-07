package protocol

/*
* https://www.home-assistant.io/docs/mqtt/discovery/
 */

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/home2mqtt/hass"
)

type IDiscoverable interface {
	MQTTProtocol
	GetObjectId() string
	GetDiscoveryInfo(uniqueID string, device *hass.Device) hass.IConfig
}

type DiscoverableNode struct {
	DiscoveryPrefix string `ini:"discovery_prefix"`
	NodeID          string `ini:"nodeid"`
}

type BasicDeviceConfig struct {
	Name     string `ini:"name,omitempty"`
	ObjectId string
}

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func PublisDiscoveryMessage(client mqtt.Client, node *DiscoverableNode, component IDiscoverable, device *hass.Device) error {
	uniqueid := node.NodeID + "_" + component.GetObjectId()
	di := component.GetDiscoveryInfo(uniqueid, device)
	c, err := json.Marshal(di)
	if err != nil {
		return err
	}
	topic := node.DiscoveryPrefix + "/" + di.GetComponent() + "/" + node.NodeID + "/" + component.GetObjectId() + "/config"
	client.Publish(topic, 0, false, c)
	return nil
}

func (device *BasicDeviceConfig) GetObjectId() string {
	return device.ObjectId
}
