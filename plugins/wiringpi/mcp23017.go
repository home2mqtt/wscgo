package main

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

func (*mcp23017addon) Configure(c interface{}) (IoImpl, error) {
	conf, err := c.(*Mcp23017Config)
	if err != nil {
		return nil, err
	}

	err = wiringpiio.Mcp23017Setup(conf)
	if err != nil {
		return nil, err
	}

	return &pinRange{
		WiringPiIO: wiringpiio,
		start: c.ExpansionBase,
		count: 16
	}
}

