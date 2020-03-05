package devices

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type ShutterStateListener func(int)

type IShutter interface {
	Device
	Open()
	Close()
	MoveBy(int)
	Stop()
	GetRange() int
	AddListener(ShutterStateListener)
}

type ShutterConfig struct {
	UpPin         string `ini:"uppin"`
	DownPin       string `ini:"downpin"`
	DirSwitchWait int    `ini:"dirswitchwait"`
	Range         int    `ini:"range"`
}

type shutter struct {
	upPin   gpio.PinOut
	downPin gpio.PinOut
	config  *ShutterConfig
	Cmd     int
	Wait    int

	Current int
	Prev    int

	PrevDir     int
	firstCmd    bool
	stopCounter int
	shouldWait  bool

	listeners []ShutterStateListener
}

func CreateShutter(config *ShutterConfig) IShutter {
	return &shutter{
		upPin:       gpioreg.ByName(config.UpPin),
		downPin:     gpioreg.ByName(config.DownPin),
		config:      config,
		Current:     0,
		Prev:        0,
		PrevDir:     0,
		firstCmd:    true,
		stopCounter: 0,
		shouldWait:  false,
	}
}

func (shutter *shutter) fireEvent() {
	for _, listener := range shutter.listeners {
		go listener(shutter.Current)
	}
}

func (shutter *shutter) up() {
	shutter.downPin.Out(gpio.Low) // turn off down
	shutter.upPin.Out(gpio.High)  // turn on up
}

func (shutter *shutter) down() {
	shutter.upPin.Out(gpio.Low)    // turn off up
	shutter.downPin.Out(gpio.High) // turn on down
}

func (shutter *shutter) stop() {
	shutter.upPin.Out(gpio.Low)   // turn off up
	shutter.downPin.Out(gpio.Low) // turn on down
}

func (shutter *shutter) Initialize() {
	shutter.stop()
	shutter.Prev = -1
}

func (shutter *shutter) setCmd(steps int) {
	if steps == 0 {
		//stop
		shutter.Cmd = 0
		if shutter.PrevDir != 0 && !shutter.firstCmd {
			shutter.Wait = shutter.config.DirSwitchWait
		} else {
			//shutter.PrevDir = 0
		}
	} else if steps > 0 {
		//up
		if shutter.Cmd < 0 {
			// direction change
			shutter.Cmd = steps
			shutter.Wait = shutter.config.DirSwitchWait
		} else {
			shutter.Cmd += steps
		}

		if (shutter.PrevDir != 1 && shutter.PrevDir != 0 && !shutter.firstCmd) || (shutter.shouldWait && shutter.PrevDir != 1) {
			shutter.Wait = shutter.config.DirSwitchWait - shutter.stopCounter
		}

		shutter.PrevDir = 1
	} else {
		//down
		if shutter.Cmd > 0 {
			//direction change
			shutter.Cmd = steps
			shutter.Wait = shutter.config.DirSwitchWait
		} else {
			shutter.Cmd += steps
		}

		if (shutter.PrevDir != -1 && shutter.PrevDir != 0 && !shutter.firstCmd) || (shutter.shouldWait && shutter.PrevDir != -1) {
			shutter.Wait = shutter.config.DirSwitchWait - shutter.stopCounter
		}

		shutter.PrevDir = -1
	}

	shutter.firstCmd = false
}

func (shutter *shutter) Tick() {
	if shutter.Wait > 0 {
		shutter.stop()
		shutter.Wait--
		if shutter.Wait == 0 {
			shutter.PrevDir = 0
		}
	} else if shutter.Cmd == 0 {
		shutter.stop()

		if shutter.stopCounter <= shutter.config.DirSwitchWait && shutter.shouldWait {
			shutter.stopCounter++
		}
		if shutter.stopCounter >= shutter.config.DirSwitchWait {
			shutter.shouldWait = false
		} else {
			shutter.shouldWait = true
		}

	} else {
		shutter.stopCounter = 0
		if shutter.Cmd > 0 {
			shutter.up()
			shutter.Cmd--
			shutter.Current++
			if shutter.Current > shutter.config.Range {
				shutter.Current = shutter.config.Range
			}
			shutter.PrevDir = 1
		} else {
			shutter.down()
			shutter.Cmd++
			shutter.Current--
			if shutter.Current < 0 {
				shutter.Current = 0
			}
			shutter.PrevDir = -1
		}

		shutter.fireEvent()
	}
}

func (shutter *shutter) AddListener(listener ShutterStateListener) {
	shutter.listeners = append(shutter.listeners, listener)
}

func (shutter *shutter) Open() {
	shutter.setCmd(shutter.config.Range)
}
func (shutter *shutter) Close() {
	shutter.setCmd(-shutter.config.Range)
}
func (shutter *shutter) MoveBy(steps int) {
	shutter.setCmd(steps)
}
func (shutter *shutter) Stop() {
	shutter.setCmd(0)
}
func (shutter *shutter) GetRange() int {
	return shutter.config.Range
}
