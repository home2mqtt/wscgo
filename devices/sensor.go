package devices

// ISensorValueListener is the signature to the callback called on measurements
type ISensorValueListener func(float64)

// ISensor is the interface of a single value sensor
type ISensor interface {
	// AddListener registers a measurement callback
	AddListener(ISensorValueListener)
	// Unit returns the human-readable representation of the measurement unit
	Unit() string
}

// BaseSensor is a default implementation of the ISensor interface
type BaseSensor struct {
	listeners []ISensorValueListener
	VUnit     string
}

func (bs *BaseSensor) Unit() string {
	return bs.VUnit
}

func (bs *BaseSensor) AddListener(listener ISensorValueListener) {
	bs.listeners = append(bs.listeners, listener)
}

func (bs *BaseSensor) NotifyListeners(value float64) {
	for _, l := range bs.listeners {
		go l(value)
	}
}
