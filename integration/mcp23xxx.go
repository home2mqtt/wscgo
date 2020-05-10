package integration

import (
	"log"

	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/periph/mcp23xxx"
	"periph.io/x/periph/conn/i2c/i2creg"
)

type mcp23xxxConfigParser struct{}

type Mcp23xxxConfig struct {
	Address int `ini:"address"`
}

func (*mcp23xxxConfigParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	c := &Mcp23xxxConfig{}
	err := section.FillData(c)
	if err != nil {
		return err
	}
	context.AddConfigInitializer(func() error {
		bus, err := i2creg.Open("")
		if err != nil {
			return err
		}
		_, err = mcp23xxx.NewI2C(bus, mcp23xxx.MCP23017, uint16(c.Address))
		if err != nil {
			return err
		}
		log.Printf("Configured mcp23017 at 0x%x", c.Address)
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("mcp23017", &mcp23xxxConfigParser{})
}
