//go:build linux
// +build linux

package cometblue

import (
	"sync"

	"github.com/home2mqtt/wscgo/devices"
)

// recoverDuration time to wait before retrying in case of communication error
const recoverDuration int = 300

var cbrange devices.ThermostatRange = devices.ThermostatRange{
	Min:  7.0,
	Max:  28.5,
	Step: 0.5,
}

// Config holds configuration values of CometBlue device
type Config struct {
	Mac      string `ini:"mac"`
	Duration int    `ini:"duration"`
}

type device struct {
	config      *Config
	dev         *Client
	counter     int
	temperature devices.BaseSensor
	targettemp  devices.BaseSensor
	target      float32
	targetset   bool
	lock        sync.Mutex
}

// Create CometBlue thermostate interface
func Create(config *Config) devices.IThermostat {
	return &device{
		config:  config,
		dev:     nil,
		counter: 0,
		temperature: devices.BaseSensor{
			VUnit: "°C",
		},
		targettemp: devices.BaseSensor{
			VUnit: "°C",
		},
		targetset: false,
	}
}

func (d *device) connect() error {
	if d.dev == nil {
		dev, err := Dial(d.config.Mac)
		if err != nil {
			return err
		}
		err = dev.Authenticate()
		if err != nil {
			dev.Close()
			return err
		}
		d.dev = dev
	}
	return nil
}

func (d *device) Initialize() error {
	d.counter = d.config.Duration
	return d.connect()
}

func (d *device) communicationError() {
	d.disconnect()
	d.counter = recoverDuration
}

func (d *device) disconnect() {
	d.dev.Close()
	d.dev = nil
}

func (d *device) Tick() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.counter > 0 {
		d.counter--
	} else {
		err := d.connect()
		if err != nil {
			d.counter = recoverDuration
			return err
		}

		t, err := d.dev.ReadTemperatures()
		if err != nil {
			d.communicationError()
			return err
		}

		// Write target temperature to device
		if d.targetset {
			err = d.dev.WriteTargetTemperature(d.target)
			if err != nil {
				d.communicationError()
				return err
			}
			t.Target = d.target
			d.targetset = false
		}

		d.processTemperature(t)

		d.disconnect()
		d.counter = d.config.Duration
	}
	return nil
}

func (d *device) processTemperature(t Temperatures) {
	// Publish results
	d.temperature.NotifyListeners(float64(t.Current))
	d.targettemp.NotifyListeners(float64(t.Target))
}

func (d *device) TargetTemperature() devices.ISensor {
	return &d.targettemp
}

func (d *device) Temperature() devices.ISensor {
	return &d.temperature
}

func (d *device) SetTargetTemperature(value float64) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.target = float32(value)
	d.targetset = true
	d.counter = 0
}

func (d *device) TemperatureRange() devices.ThermostatRange {
	return cbrange
}
