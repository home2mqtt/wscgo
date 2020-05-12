package integration

import (
	"log"

	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/periph/pca9685"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
)

type pca9685ConfigParser struct{}

type Pca9685Config struct {
	Address   int `ini:"address"`
	Frequency int `ini:"frequency"`
}

func (*pca9685ConfigParser) ParseConfiguration(section config.ConfigurationSection, context config.ConfigurationContext) error {
	c := &Pca9685Config{}
	err := section.FillData(c)
	if err != nil {
		return err
	}
	context.AddConfigInitializer(func() error {
		bus, err := i2creg.Open("")
		if err != nil {
			return err
		}
		dev, err := pca9685.NewI2C(bus, uint16(c.Address))
		if err != nil {
			return err
		}
		err = dev.SetPwmFreq(physic.Frequency(c.Frequency))
		if err != nil {
			return err
		}
		log.Printf("Configured pca9685 at 0x%x", c.Address)
		return dev.RegisterPins()
	})
	return nil
}

func init() {
	config.RegisterConfigurationPartParser("pca9685", &pca9685ConfigParser{})
}
