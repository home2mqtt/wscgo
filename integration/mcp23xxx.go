package integration

import (
	"fmt"
	"log"
	"strings"

	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/periph/mcp23xxx"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

type mcp23xxxConfigParser struct {
	variant mcp23xxx.Variant
}

type Mcp23xxxConfig struct {
	Address int `ini:"address"`
}

func (p *mcp23xxxConfigParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	c := &Mcp23xxxConfig{}
	err := section.FillData(c)
	if err != nil {
		return err
	}
	switch p.variant {
	case mcp23xxx.MCP23008, mcp23xxx.MCP23009, mcp23xxx.MCP23016, mcp23xxx.MCP23017, mcp23xxx.MCP23018:
		context.AddConfigInitializer(func() error {
			bus, err := i2creg.Open("")
			if err != nil {
				return err
			}
			_, err = mcp23xxx.NewI2C(bus, p.variant, uint16(c.Address))
			if err != nil {
				return err
			}
			log.Printf("Configured %s at 0x%x", p.variant, c.Address)
			return nil
		})
		return nil
	case mcp23xxx.MCP23S08, mcp23xxx.MCP23S09, mcp23xxx.MCP23S17, mcp23xxx.MCP23S18:
		context.AddConfigInitializer(func() error {
			bus, err := spireg.Open("")
			if err != nil {
				return err
			}

			c, err := bus.Connect(physic.MegaHertz, spi.Mode3, 8)
			if err != nil {
				return err
			}
			_, err = mcp23xxx.NewSPI(c, p.variant)
			if err != nil {
				return err
			}
			log.Printf("Configured %s", p.variant)
			return nil
		})
		return nil
	default:
		return fmt.Errorf("Unknown MCP23 variant: %s", p.variant)
	}
}

func register(variant mcp23xxx.Variant) {
	err := config.RegisterConfigurationPartParser(strings.ToLower(string(variant)), &mcp23xxxConfigParser{
		variant: variant,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	register(mcp23xxx.MCP23008)
	register(mcp23xxx.MCP23009)
	register(mcp23xxx.MCP23016)
	register(mcp23xxx.MCP23017)
	register(mcp23xxx.MCP23018)
	register(mcp23xxx.MCP23S08)
	register(mcp23xxx.MCP23S09)
	register(mcp23xxx.MCP23S17)
	register(mcp23xxx.MCP23S18)
}
