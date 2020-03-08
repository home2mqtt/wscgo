package wiringpi

import "gitlab.com/grill-tamasi/wscgo/config"

type Pca9685Config struct {
	Address       int     `ini:"address"`
	ExpansionBase int     `ini:"expansionBase"`
	Frequency     float32 `ini:"frequency"`
}

type pca9685configPartParser struct{}

func (*pca9685configPartParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	c := &Pca9685Config{}
	section.FillData(c)
	context.AddConfigInitializer(func() error {
		Pca9685Setup(c)
		return nil
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("mcp23017", &pca9685configPartParser{})
}