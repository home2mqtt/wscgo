package homeassistant

/*
* https://www.home-assistant.io/docs/mqtt/discovery/
 */

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.com/grill-tamasi/wscgo/protocol"
)

type IDiscoverable interface {
	protocol.MQTTProtocol
	GetComponent() string
	GetObjectId() string
	GetDiscoveryInfo() interface{}
	GetJsonDiscoveryInfo() ([]byte, error)
}

type DiscoverableNode struct {
	DiscoveryPrefix string
	NodeID          string
}

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func PublisDiscoveryMessage(client mqtt.Client, node *DiscoverableNode, component IDiscoverable) error {
	c, err := component.GetJsonDiscoveryInfo()
	if err != nil {
		return err
	}
	topic := node.DiscoveryPrefix + "/" + component.GetComponent() + "/" + node.NodeID + "/" + component.GetObjectId() + "/config"
	client.Publish(topic, 0, true, c)
	return nil
}
