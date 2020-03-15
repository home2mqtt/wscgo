package main

import (
	"errors"

	"gitlab.com/grill-tamasi/wscgo/plugins"
)

type Pca9685Config struct {
	Address       int     `ini:"address"`
	ExpansionBase int     `ini:"expansionBase"`
	Frequency     float32 `ini:"frequency"`
}

type pca9685addon struct{}

func (*pca9685addon) GetIdentifier() string {
	return "pca9685"
}

func (*pca9685addon) CreateConfigStruct() interface{} {
	return &Pca9685Config{}
}

func (*pca9685addon) Configure(c interface{}) (plugins.IoImpl, error) {
	conf, ok := c.(*Pca9685Config)
	if !ok {
		return nil, errors.New("PCA9685: invalid config type")
	}

	err := Pca9685Setup(conf)
	if err != nil {
		return nil, err
	}

	return &pinRange{
		WiringPiIO: wiringpiio,
		start:      conf.ExpansionBase,
		count:      16,
	}, nil
}
