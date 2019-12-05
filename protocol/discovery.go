package protocol

/*
* https://www.home-assistant.io/docs/mqtt/discovery/
 */

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IDiscoverable interface {
	MQTTProtocol
	GetComponent() string
	GetObjectId() string
	GetDiscoveryInfo() interface{}
}

type DiscoverableNode struct {
	DiscoveryPrefix string
	NodeID          string
}

type BasicDeviceConfig struct {
	Name     string `ini:"name,omitempty"`
	ObjectId string
}

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func PublisDiscoveryMessage(client mqtt.Client, node *DiscoverableNode, component IDiscoverable) error {
	di := component.GetDiscoveryInfo()
	c, err := json.Marshal(di)
	if err != nil {
		return err
	}
	topic := node.DiscoveryPrefix + "/" + component.GetComponent() + "/" + node.NodeID + "/" + component.GetObjectId() + "/config"
	client.Publish(topic, 0, false, c)
	return nil
}

func (device *BasicDeviceConfig) GetObjectId() string {
	return device.ObjectId
}
