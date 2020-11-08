// +build linux

package cometblue

import (
	"log"
	"sync"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

// CometBlueContext holds static information to control CB devices
type CometBlueContext struct {
	once sync.Once

	ThermostatService ble.UUID
	TemperaturesChar  ble.UUID
	BatteryChar       ble.UUID
	PinChar           ble.UUID
}

var cbContextInstance CometBlueContext

func (c *CometBlueContext) init() error {
	var err error
	c.once.Do(func() {
		log.Println("Initializing bluetooth")
		d, err := dev.NewDevice("default")
		if err != nil {
			return
		}
		ble.SetDefaultDevice(d)

		c.ThermostatService, err = ble.Parse("47e9ee00-47e9-11e4-8939-164230d1df67")
		if err != nil {
			return
		}
		c.TemperaturesChar, err = ble.Parse("47e9ee2b-47e9-11e4-8939-164230d1df67")
		if err != nil {
			return
		}
		c.BatteryChar, err = ble.Parse("47e9ee2c-47e9-11e4-8939-164230d1df67")
		if err != nil {
			return
		}
		c.PinChar, err = ble.Parse("47e9ee30-47e9-11e4-8939-164230d1df67")
		if err != nil {
			return
		}
	})
	return err
}

// GetCBContext retrieve singleton CometBlueContext instance and makes sure that Bluetooth is initialized
func GetCBContext() *CometBlueContext {
	err := cbContextInstance.init()
	if err != nil {
		log.Fatalf("Failed to Initialize Bluetooth: %s", err)
	}
	return &cbContextInstance
}
