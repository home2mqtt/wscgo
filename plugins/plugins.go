package plugins

type AddonsGetter func() []Addon

const AddonsGetterName = "GetAddons"

type IoImpl interface {
	DigitalWrite(pin int, value bool)
	DigitalRead(pin int) bool
	PinMode(pin int, mode int) error
	PwmWrite(pin int, value int) error
	PinRange() (int, int)
	PwmResolution() int
}

type Addon interface {
	GetIdentifier() string
	CreateConfigStruct() interface{}
	Configure(interface{}) (IoImpl, error)
}
