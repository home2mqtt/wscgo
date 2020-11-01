// +build linux

package cometblue

import (
	"github.com/balazsgrill/wscgo/devices"
)

// recoverDuration time to wait before retrying in case of communication error
const recoverDuration int = 300

var cbrange devices.ThermostatRange = devices.ThermostatRange{
	Min:  7.0,
	Max:  23.0,
	Step: 0.5,
}

type CometblueConfig struct {
	Mac      string `ini:"mac"`
	Duration int    `ini:"duration"`
}

type blueCometDevice struct {
	config      *CometblueConfig
	dev         *CometblueClient
	counter     int
	temperature devices.BaseSensor
	targettemp  devices.BaseSensor
	target      float32
	targetset   bool
}

func CreateCometblue(config *CometblueConfig) devices.IThermostat {
	return &blueCometDevice{
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

func (d *blueCometDevice) connect() error {
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

func (d *blueCometDevice) Initialize() error {
	d.counter = d.config.Duration
	return d.connect()
}

func (d *blueCometDevice) communicationError() {
	d.dev.Close()
	d.dev = nil
	d.counter = recoverDuration
}

func (d *blueCometDevice) Tick() error {
	if d.counter > 0 {
		d.counter--
	} else {
		err := d.connect()
		if err != nil {
			d.counter = recoverDuration
			return err
		}
		if d.targetset {
			err = d.dev.WriteTargetTemperature(d.target)
			if err != nil {
				d.communicationError()
				return err
			}
			d.targetset = false
		}

		t, err := d.dev.ReadTemperatures()
		if err != nil {
			d.communicationError()
			return err
		}
		d.counter = d.config.Duration
		d.temperature.NotifyListeners(float64(t.Current))
		d.temperature.NotifyListeners(float64(t.Target))
	}
	return nil
}

func (d *blueCometDevice) TargetTemperature() devices.ISensor {
	return &d.targettemp
}

func (d *blueCometDevice) Temperature() devices.ISensor {
	return &d.temperature
}

func (d *blueCometDevice) SetTargetTemperature(value float64) {
	d.target = float32(value)
	d.targetset = true
}

func (d *blueCometDevice) TemperatureRange() devices.ThermostatRange {
	return cbrange
}
