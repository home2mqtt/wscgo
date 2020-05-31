package devices

import "fmt"

// Device is the interface of a periodically operated peripheral controlled by wscgo.
type Device interface {
	Tick() error
	Initialize() error
}

func invalidPinError(pinName string) error {
	return fmt.Errorf("Pin not found: %s", pinName)
}
