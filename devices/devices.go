package devices

type IoContext interface {
	DigitalWrite(pin int, value bool)
	DigitalRead(pin int) bool
	PinMode(pin int, mode int)
}

type Device interface {
	Tick()
	Initialize()
}
