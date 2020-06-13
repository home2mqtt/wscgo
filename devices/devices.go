package devices

import "fmt"

// Device is the interface of a periodically operated peripheral controlled by wscgo.
type Device interface {
	// Tick is called periodically
	Tick() error
	// Initialize is called upon startup
	Initialize() error
}

func invalidPinError(pinName string) error {
	return fmt.Errorf("Pin not found: %s", pinName)
}
