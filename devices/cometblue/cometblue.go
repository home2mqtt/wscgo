// +build linux

package cometblue

import (
	"context"
	"math"
	"time"

	"errors"

	"github.com/go-ble/ble"
)

const timeout time.Duration = 5 * time.Second
const unusedbyte byte = 0x80

// Client provides API to control a CometBlue device
type Client struct {
	client      ble.Client
	pin         *ble.Characteristic
	battery     *ble.Characteristic
	temperature *ble.Characteristic
}

// Temperatures data struct for a temperature measurement
type Temperatures struct {
	Current float32
	Target  float32
}

// Dial attempts to connect to a CometBlue device
func Dial(address string) (*Client, error) {
	bt := GetCBContext()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := ble.Dial(ctx, ble.NewAddr(address))
	if err != nil {
		return nil, err
	}
	var ctemp *ble.Characteristic
	var cbatt *ble.Characteristic
	var cpin *ble.Characteristic
	p, err := client.DiscoverProfile(true)
	if err != nil {
		return nil, err
	}
	if ctemp = p.FindCharacteristic(&ble.Characteristic{UUID: bt.TemperaturesChar}); ctemp == nil {
		client.CancelConnection()
		return nil, errors.New("Couldn't find Temperature characteristic")
	}
	if cbatt = p.FindCharacteristic(&ble.Characteristic{UUID: bt.BatteryChar}); cbatt == nil {
		client.CancelConnection()
		return nil, errors.New("Couldn't find Battery characteristic")
	}
	if cpin = p.FindCharacteristic(&ble.Characteristic{UUID: bt.PinChar}); cpin == nil {
		client.CancelConnection()
		return nil, errors.New("Couldn't find PIN characteristic")
	}

	c := &Client{
		client:      client,
		pin:         cpin,
		battery:     cbatt,
		temperature: ctemp,
	}
	return c, nil
}

// Authenticate writes in default PIN
func (client *Client) Authenticate() error {
	return client.client.WriteCharacteristic(client.pin, []byte{0, 0, 0, 0}, false)
}

// ReadBatteryRaw raw battery characteristic
func (client *Client) ReadBatteryRaw() ([]byte, error) {
	return client.client.ReadCharacteristic(client.battery)
}

// ReadTemperaturesRaw raw temperature characteristic
func (client *Client) ReadTemperaturesRaw() ([]byte, error) {
	return client.client.ReadCharacteristic(client.temperature)
}

// ReadTemperatures temperature
func (client *Client) ReadTemperatures() (Temperatures, error) {
	data, err := client.ReadTemperaturesRaw()
	if err != nil {
		return Temperatures{}, err
	}
	if data[0] == unusedbyte {
		return Temperatures{}, errors.New("Failed to read temperature from device")
	}
	return Temperatures{
		Current: float32(data[0]) * 0.5,
		Target:  float32(data[1]) * 0.5,
	}, err
}

// WriteTargetTemperature sets target temperature
func (client *Client) WriteTargetTemperature(target float32) error {
	// TODO range-check target
	targetValue := byte(0xff & int32(math.Round(float64(target*2))))
	data := []byte{
		unusedbyte,
		targetValue,
		targetValue,
		targetValue,
		0,
		unusedbyte,
		unusedbyte,
	}
	return client.client.WriteCharacteristic(client.temperature, data, false)
}

// Close connection
func (client *Client) Close() {
	client.client.CancelConnection()
}
