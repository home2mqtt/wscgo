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

type baseSensor struct {
	listeners []ISensorValueListener
	unit      string
}

func (bs *baseSensor) Unit() string {
	return bs.unit
}

func (bs *baseSensor) AddListener(listener ISensorValueListener) {
	bs.listeners = append(bs.listeners, listener)
}

func (bs *baseSensor) notifyListeners(value float64) {
	for _, l := range bs.listeners {
		l(value)
	}
}
