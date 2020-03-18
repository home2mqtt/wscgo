package devices

import (
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

var i2cbusRegister = make(map[string]i2c.Bus)

func getI2CBus(name string) (i2c.Bus, error) {
	bus, ok := i2cbusRegister[name]
	if !ok {
		bus, err := i2creg.Open(name)
		if err != nil {
			return nil, err
		}
		i2cbusRegister[name] = bus
	}
	return bus, nil
}
