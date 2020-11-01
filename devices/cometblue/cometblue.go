// +build linux

package cometblue

import (
	"context"
	"log"
	"math"
	"time"

	"errors"

	"github.com/go-ble/ble"
)

const timeout time.Duration = 5 * time.Second

type CometblueClient struct {
	client      ble.Client
	service     *ble.Service
	pin         *ble.Characteristic
	battery     *ble.Characteristic
	temperature *ble.Characteristic
}

type Temperatures struct {
	Current float32
	Target  float32
}

func Dial(address string) (*CometblueClient, error) {
	bt := GetCBContext()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := ble.Dial(ctx, ble.NewAddr(address))
	if err != nil {
		return nil, err
	}
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
	for _, char := range characteristics {
		log.Println(char.UUID)
	}
	if len(characteristics) != 3 {
		client.CancelConnection()
		return nil, errors.New("Couldn't read Characteristics")
	}
	return &CometblueClient{
		client:      client,
		service:     services[0],
		pin:         characteristics[2],
		battery:     characteristics[1],
		temperature: characteristics[0],
	}, nil
}

func (client *CometblueClient) Authenticate() error {
	return client.client.WriteCharacteristic(client.pin, []byte{0, 0, 0, 0}, false)
}

func (client *CometblueClient) ReadBatteryRaw() ([]byte, error) {
	return client.client.ReadCharacteristic(client.battery)
}

func (client *CometblueClient) ReadTemperaturesRaw() ([]byte, error) {
	return client.client.ReadCharacteristic(client.temperature)
}

func (client *CometblueClient) ReadTemperatures() (Temperatures, error) {
	data, err := client.ReadTemperaturesRaw()
	if err != nil {
		return Temperatures{}, err
	}
	return Temperatures{
		Current: float32(data[0]) * 0.5,
		Target:  float32(data[1]) * 0.5,
	}, err
}

func (client *CometblueClient) WriteTargetTemperature(target float32) error {
	// TODO range-check target
	targetValue := byte(0xff & int32(math.Round(float64(target*2))))
	data := []byte{
		0x80,
		0x80,
		targetValue,
		targetValue,
		0,
		0x80,
		0x80,
	}
	return client.client.WriteCharacteristic(client.temperature, data, false)
}

func (client *CometblueClient) Close() {
	client.client.CancelConnection()
}
