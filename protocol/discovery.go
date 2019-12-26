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
	GetDiscoveryInfo(uniqueID string, device *DeviceDiscoveryInfo) interface{}
}

type DiscoverableNode struct {
	DiscoveryPrefix string `ini:"discovery_prefix"`
	NodeID          string `ini:"nodeid"`
}

type BasicDeviceConfig struct {
	Name     string `ini:"name,omitempty"`
	ObjectId string
}

type DeviceDiscoveryInfo struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Connections  []string `json:"connections,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
	SwVersion    string   `json:"sw_version,omitempty"`
}

type BasicDiscoveryInfo struct {
	Device   *DeviceDiscoveryInfo `json:"device,omitempty"`
	UniqueID string               `json:"unique_id,omitempty"`
}

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func PublisDiscoveryMessage(client mqtt.Client, node *DiscoverableNode, component IDiscoverable, device *DeviceDiscoveryInfo) error {
	uniqueid := node.NodeID + "_" + component.GetComponent() + "_" + component.GetObjectId()
	di := component.GetDiscoveryInfo(uniqueid, device)
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
