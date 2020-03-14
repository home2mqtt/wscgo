package wiringpi

import (
	"gitlab.com/grill-tamasi/wscgo/config"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type Mcp23017Config struct {
	Address       int `ini:"address"`
	ExpansionBase int `ini:"expansionBase"`
}

type mcp23017configPartParser struct{}

func (*mcp23017configPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	c := &Mcp23017Config{}
	section.FillData(c)
	context.AddConfigInitializer(func() error {
		Mcp23017Setup(c)
		for i := 0; i < 16; i++ {
			gpioreg.Register(New(c.ExpansionBase + i))
		}
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("mcp23017", &mcp23017configPartParser{})
}
