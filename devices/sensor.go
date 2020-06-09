package devices

type ISensorValueListener func(float64)

type ISensor interface {
	AddListener(ISensorValueListener)
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
