package protocol

import (
	"encoding/json"
	"testing"

	"gitlab.com/grill-tamasi/wscgo/devices"
	"gitlab.com/grill-tamasi/wscgo/tests"
)

type testshutter struct {
	devices.IoContext
}

func (*testshutter) Tick()                                      {}
func (*testshutter) Initialize()                                {}
func (*testshutter) Open()                                      {}
func (*testshutter) Close()                                     {}
func (*testshutter) OpenBy(int)                                 {}
func (*testshutter) CloseBy(int)                                {}
func (*testshutter) Stop()                                      {}
func (*testshutter) GetRange() int                              { return 10 }
func (*testshutter) AddListener(l devices.ShutterStateListener) {}

func TestDiscoveryJson(t *testing.T) {
	conf := &CoverConfig{
		CommandTopic:  "test/command",
		Name:          "name",
		PositionTopic: "test/pos",
		ObjectId:      "objectId",
	}
	c := IntegrateCover(&testshutter{
		IoContext: tests.CreateTestIo(2),
	}, conf)
	data, err := c.GetJsonDiscoveryInfo()
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
		t.Errorf("Command topic is invalid!")
	}
}
