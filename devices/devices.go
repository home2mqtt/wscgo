package devices

type Device interface {
	Tick()
	Initialize()
}
