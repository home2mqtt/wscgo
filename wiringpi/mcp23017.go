package wiringpi

type Mcp23017Config struct {
	Address       int `ini:"address"`
	ExpansionBase int `ini:"expansionBase"`
}
