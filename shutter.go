package main

import "log"

type shutterCallback func(int)

type shutter struct {
	id string
	ioContext
	UpPin         int
	DownPin       int
	DirSwitchWait int

	Cmd  int
	Wait int

	Range   int
	Current int
	Prev    int

	topic      string
	groupTopic string
	Callback   shutterCallback

	PrevDir int

	firstCmd    bool
	stopCounter int
	shouldWait  bool
}

func (shutter *shutter) up() {
	shutter.digitalWrite(shutter.DownPin, LOW) // turn off down
	shutter.digitalWrite(shutter.UpPin, HIGH)  // turn on up
}

func (shutter *shutter) down() {
	shutter.digitalWrite(shutter.UpPin, LOW)    // turn off up
	shutter.digitalWrite(shutter.DownPin, HIGH) // turn on down
}

func (shutter *shutter) stop() {
	shutter.digitalWrite(shutter.UpPin, LOW)   // turn off up
	shutter.digitalWrite(shutter.DownPin, LOW) // turn on down
}

func (shutter *shutter) init() {
	shutter.pinMode(shutter.UpPin, OUTPUT)
	shutter.pinMode(shutter.DownPin, OUTPUT)
	shutter.Prev = -1
}

func (shutter *shutter) setCmd(steps int) {
	if steps == 0 {
		//stop
		shutter.Cmd = 0
		if shutter.PrevDir != 0 && !shutter.firstCmd {
			shutter.Wait = shutter.DirSwitchWait
		} else {
			//shutter.PrevDir = 0
		}
	} else if steps > 0 {
		//up
		if shutter.Cmd < 0 {
			// direction change
			shutter.Cmd = steps
			shutter.Wait = shutter.DirSwitchWait
		} else {
			shutter.Cmd += steps
		}

		if (shutter.PrevDir != 1 && shutter.PrevDir != 0 && !shutter.firstCmd) || (shutter.shouldWait && shutter.PrevDir != 1) {
			shutter.Wait = shutter.DirSwitchWait - shutter.stopCounter
		}

		shutter.PrevDir = 1
	} else {
		//down
		if shutter.Cmd > 0 {
			//direction change
			shutter.Cmd = steps
			shutter.Wait = shutter.DirSwitchWait
		} else {
			shutter.Cmd += steps
		}

		if (shutter.PrevDir != -1 && shutter.PrevDir != 0 && !shutter.firstCmd) || (shutter.shouldWait && shutter.PrevDir != -1) {
			shutter.Wait = shutter.DirSwitchWait - shutter.stopCounter
		}

		shutter.PrevDir = -1
	}

	shutter.firstCmd = false
}

func (shutter *shutter) tick() {
	if shutter.Wait > 0 {
		shutter.stop()
		shutter.Wait--
		if shutter.Wait == 0 {
			shutter.PrevDir = 0
		}
		log.Println("WAIT")
	} else if shutter.Cmd == 0 {
		shutter.stop()

		if shutter.stopCounter <= shutter.DirSwitchWait && shutter.shouldWait {
			shutter.stopCounter++
		}
		if shutter.stopCounter >= shutter.DirSwitchWait {
			shutter.shouldWait = false
			log.Println("STOP")
		} else {
			shutter.shouldWait = true
			log.Println("STOP - Should wait")
		}

	} else {
		shutter.stopCounter = 0
		if shutter.Cmd > 0 {
			shutter.up()
			shutter.Cmd--
			shutter.Current++
			if shutter.Current > shutter.Range {
				shutter.Current = shutter.Range
			}
			shutter.PrevDir = 1
			log.Println("UP")
		} else {
			shutter.down()
			shutter.Cmd++
			shutter.Current--
			if shutter.Current < 0 {
				shutter.Current = 0
			}
			shutter.PrevDir = -1
			log.Println("DOWN")
		}

		if shutter.Callback != nil {
			if shutter.Prev != shutter.Current {
				shutter.Prev = shutter.Current
				shutter.Callback(shutter.Current)
			}
		}
	}
}
