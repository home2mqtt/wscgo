package devices

import (
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/experimental/devices/ccs811"
)

// ICCS811 denotes a CCS811 environmental sensor
type ICCS811 interface {
	Device
	// ECO2 returns the sensor interface to the "equivalent CO2" measurement value
	ECO2() ISensor
	// VOC returns the sensor interface to the "Total Volatile Organic Compound" measuremenr value
	VOC() ISensor
}

// CCS811Config holds the configuration data of a sensor
type CCS811Config struct {
	Address  int    `ini:"address"`
	Bus      string `ini:"i2cbus"`
	Duration int    `ini:"duration"`
}

type ccs811Device struct {
	dev          *ccs811.Dev
	config       *CCS811Config
	eco2         BaseSensor
	voc          BaseSensor
	measureCount int
}

// CreateCCS811 configures connection the sensor
func CreateCCS811(config *CCS811Config) (ICCS811, error) {
	opts := &ccs811.Opts{
		Addr:               uint16(config.Address),
		MeasurementMode:    ccs811.MeasurementModeConstant1000,
		InterruptWhenReady: false,
		UseThreshold:       false,
	}

	bus, err := i2creg.Open(config.Bus)
	if err != nil {
		return nil, err
	}
	dev, err := ccs811.New(bus, opts)
	if err != nil {
		return nil, err
	}

	return &ccs811Device{
		dev:          dev,
		config:       config,
		measureCount: config.Duration,
		eco2: BaseSensor{
			VUnit: "ppm",
		},
		voc: BaseSensor{
			VUnit: "ppb",
		},
	}, nil
}

func (dev *ccs811Device) Tick() error {
	values := &ccs811.SensorValues{}
	if dev.measureCount == 0 {
		dev.measureCount = dev.config.Duration
		err := dev.dev.Sense(values)
		if err != nil {
			return err
		}
		(&dev.eco2).NotifyListeners(float64(values.ECO2))
		(&dev.voc).NotifyListeners(float64(values.VOC))
	} else {
		dev.measureCount--
	}
	return nil
}

func (dev *ccs811Device) Initialize() error {
	dev.measureCount = dev.config.Duration
	return nil
}

func (dev *ccs811Device) ECO2() ISensor {
	return &dev.eco2
}

func (dev *ccs811Device) VOC() ISensor {
	return &dev.voc
}
