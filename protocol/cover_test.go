package protocol

import (
	"encoding/json"
	"testing"

	"gitlab.com/grill-tamasi/wscgo/devices"
)

type testshutter struct {
}

func (*testshutter) Tick() error                                { return nil }
func (*testshutter) Initialize() error                          { return nil }
func (*testshutter) Open()                                      {}
func (*testshutter) Close()                                     {}
func (*testshutter) MoveBy(int)                                 {}
func (*testshutter) Stop()                                      {}
func (*testshutter) GetRange() int                              { return 10 }
func (*testshutter) AddListener(l devices.ShutterStateListener) {}

func TestDiscoveryJson(t *testing.T) {
	conf := &CoverConfig{
		BasicDeviceConfig: BasicDeviceConfig{
			Name:     "name",
			ObjectId: "objectId",
		},
		CommandTopic:  "test/command",
		PositionTopic: "test/pos",
	}
	c := IntegrateCover(&testshutter{}, conf)
	info := c.GetDiscoveryInfo("test_"+conf.ObjectId, &DeviceDiscoveryInfo{
		Manufacturer: "wscgo",
		Model:        "wscgo",
		SwVersion:    "0.0.0-test",
		Name:         "wscgo-test",
	})
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}

	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		t.Fatal(err)
	}

	m := f.(map[string]interface{})
	if m["command_topic"] != conf.CommandTopic {
		t.Error("Command topic is invalid!")
	}
	if m["name"] != conf.Name {
		t.Error("Name is not set properly!")
	}
}
