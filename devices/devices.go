package devices

import "fmt"

type Device interface {
	Tick() error
	Initialize() error
}

func invalidPinError(pinName string) error {
	return fmt.Errorf("Pin not found: %s", pinName)
}
