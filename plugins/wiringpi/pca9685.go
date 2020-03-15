package main

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
	return &pca9685Config{}
}

func (*pca9685addon) Configure(c interface{}) (IoImpl, error) {
	conf, err := c.(*Pca9685Config)
	if err != nil {
		return nil, err
	}

	err = wiringpiio.Pca9685Setup(conf)
	if err != nil {
		return nil, err
	}

	return &pinRange{
		WiringPiIO: wiringpiio,
		start: c.ExpansionBase,
		count: 16
	}
}

