package main

import (
	"errors"
	"log"

	"github.com/grill-tamasi/wscgo/plugins"
)

type Mcp23017Config struct {
	Address       int `ini:"address"`
	ExpansionBase int `ini:"expansionBase"`
}

type mcp23017addon struct{}

func (*mcp23017addon) GetIdentifier() string {
	return "mcp23017"
}

func (*mcp23017addon) CreateConfigStruct() interface{} {
	return &Mcp23017Config{}
}

func (*mcp23017addon) Configure(c interface{}) (plugins.IoImpl, error) {
	conf, ok := c.(*Mcp23017Config)
	if !ok {
		return nil, errors.New("MCP23017: invalid config type")
	}

	err := Mcp23017Setup(conf)
	if err != nil {
		return nil, err
	}
	log.Printf("Configured MCP23017 at 0x%x\n", conf.Address)

	return &pinRange{
		WiringPiIO: wiringpiio,
		start:      conf.ExpansionBase,
		count:      16,
		pwmres:     0,
	}, nil
}
