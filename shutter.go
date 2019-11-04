package main

type shutterCallback func(int)

type shutter struct {
	ioContext
	UpPin         int
	DownPin       int
	DirSwitchWait int

	Cmd  int
	Wait int

	Range   int
	Current int
	Prev    int

	topic    string
	Callback shutterCallback

	PrevDir int
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
	} else if steps > 0 {
		//up
		if shutter.Cmd < 0 {
			// direction change
			shutter.Cmd = steps
			shutter.Wait = shutter.DirSwitchWait
		} else {
			shutter.Cmd += steps
		}

		if shutter.PrevDir < 0 {
			shutter.Wait = shutter.DirSwitchWait
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

		if shutter.PrevDir > 0 {
			shutter.Wait = shutter.DirSwitchWait
		}

		shutter.PrevDir = -1
	}
}

func (shutter *shutter) tick() {
	if shutter.Wait > 0 {
		shutter.stop()
		shutter.Wait--
	} else if shutter.Cmd == 0 {
		shutter.stop()
	} else {
		if shutter.Cmd > 0 {
			shutter.up()
			shutter.Cmd--
			shutter.Current++
			if shutter.Current > shutter.Range {
				shutter.Current = shutter.Range
			}
		} else {
			shutter.down()
			shutter.Cmd++
			shutter.Current--
			if shutter.Current < 0 {
				shutter.Current = 0
			}
		}

		if shutter.Callback != nil {
			if shutter.Prev != shutter.Current {
				shutter.Prev = shutter.Current
				shutter.Callback(shutter.Current)
			}
		}
	}
}
