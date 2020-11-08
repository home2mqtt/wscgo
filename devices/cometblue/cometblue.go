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
	service     *ble.Service
	pin         *ble.Characteristic
	battery     *ble.Characteristic
	temperature *ble.Characteristic
	Handler     TemperatureHandler
}

// Temperatures data struct for a temperature measurement
type Temperatures struct {
	Current float32
	Target  float32
}

// TemperatureHandler is called to process a temperature measurement
type TemperatureHandler func(t Temperatures)

// Dial attempts to connect to a CometBlue device
func Dial(address string) (*Client, error) {
	bt := GetCBContext()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := ble.Dial(ctx, ble.NewAddr(address))
	if err != nil {
		return nil, err
	}
	//p, err := client.DiscoverProfile(true)
	//p.
	services, err := client.DiscoverServices([]ble.UUID{bt.ThermostatService})
	if err != nil {
		client.CancelConnection()
		return nil, err
	}
	if len(services) != 1 {
		client.CancelConnection()
		return nil, errors.New("Couldn't read Service")
	}
	// Returned list of characteristics seems to be in order by UUID not by request
	characteristics, err := client.DiscoverCharacteristics([]ble.UUID{
		bt.TemperaturesChar,
		bt.BatteryChar,
		bt.PinChar,
	}, services[0])
	if err != nil {
		client.CancelConnection()
		return nil, err
	}
	if len(characteristics) != 3 {
		client.CancelConnection()
		return nil, errors.New("Couldn't read Characteristics")
	}

	c := &Client{
		client:      client,
		service:     services[0],
		pin:         characteristics[2],
		battery:     characteristics[1],
		temperature: characteristics[0],
	}
	// Discover descriptors to fill CCCD
	_, err = client.DiscoverDescriptors(nil, c.temperature)
	if err != nil {
		client.CancelConnection()
		return nil, err
	}
	err = client.Subscribe(c.temperature, true, c.temperatureNotificationHandler)
	if err != nil {
		client.CancelConnection()
		return nil, err
	}
	return c, nil
}

func (client *Client) temperatureNotificationHandler(data []byte) {
	if client.Handler != nil {
		client.Handler(Temperatures{
			Current: float32(data[0]) * 0.5,
			Target:  float32(data[1]) * 0.5,
		})
	}
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
