package devices

type IoContext interface {
	DigitalWrite(pin int, value bool)
	DigitalRead(pin int) bool
	PinMode(pin int, output bool)
}

type Device interface {
	Tick()
	Initialize()
}
