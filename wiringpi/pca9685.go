package wiringpi

import (
	"log"

	"gitlab.com/grill-tamasi/wscgo/config"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

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
		for i := 0; i < 16; i++ {
			gpioreg.Register(New(c.ExpansionBase + i))
		}
		return nil
	})
	return nil
}

func init() {
	log.Println("Pugin: PCA9685")
	config.RegisterConfigurationPartParser("pca9685", &pca9685configPartParser{})
}
