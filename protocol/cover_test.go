package protocol

import (
	"encoding/json"
	"testing"

	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/tests"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

type testshutter struct {
	wiringpi.IoContext
}

func (*testshutter) Tick()                                      {}
func (*testshutter) Initialize()                                {}
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
	c := IntegrateCover(&testshutter{
		IoContext: tests.CreateTestIo(2),
	}, conf)
	info := c.GetDiscoveryInfo()
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

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
