package devices

type ThermostatRange struct {
	Min  float64
	Max  float64
	Step float64
}

type IThermostat interface {
	Device

	TemperatureRange() ThermostatRange
	SetTargetTemperature(float64)
	TargetTemperature() ISensor
	Temperature() ISensor
}
