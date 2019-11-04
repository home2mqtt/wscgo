package main

import "testing"

type testIo struct {
	modes  []int
	values []int
}

func (io *testIo) digitalWrite(pin int, value int) {
	io.values[pin] = value
}

func (io *testIo) digitalRead(pin int) int {
	return io.values[pin]
}

func (io *testIo) pinMode(pin int, mode int) {
	io.modes[pin] = mode
}

func TestInit(t *testing.T) {
	io := testIo{
		modes:  []int{INPUT, INPUT, INPUT},
		values: []int{LOW, LOW, LOW},
	}
	s := shutter{
		ioContext:     &io,
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}

	s.init()
	if io.modes[0] != OUTPUT {
		t.Fatal()
	}
	if io.modes[1] != OUTPUT {
		t.Fatal()
	}
}

func TestUp(t *testing.T) {
	io := testIo{
		modes:  []int{INPUT, INPUT, INPUT},
		values: []int{LOW, LOW, LOW},
	}
	s := shutter{
		ioContext:     &io,
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}

	s.init()

	s.setCmd(10)
	for i := 0; i < 10; i++ {
		s.tick()
		if io.values[0] != HIGH {
			t.Errorf("up is low %d\n", i)
		}
		if io.values[1] != LOW {
			t.Errorf("down is high %d\n", i)
		}
	}
	s.tick()
	if io.values[0] != LOW {
		t.Errorf("up is high\n")
	}
}

func checkPins(msg string, t *testing.T, io testIo, up int, down int) {
	if io.values[0] != up || io.values[1] != down {
		t.Errorf("%s UP[exc-actal]: %d - %d, DOWN[exc-actal]: %d - %d\n", msg, up, io.values[0], down, io.values[1])
	}
}

func TestZero(t *testing.T) {
	io := testIo{
		modes:  []int{INPUT, INPUT, INPUT},
		values: []int{LOW, LOW, LOW},
	}
	s := shutter{
		ioContext:     &io,
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}

	s.init()

	for i := 0; i < 50; i++ {
		s.tick()
		checkPins("req zero ", t, io, LOW, LOW)
	}

}

func TestUpDown(t *testing.T) {
	io := testIo{
		modes:  []int{INPUT, INPUT, INPUT},
		values: []int{LOW, LOW, LOW},
	}
	s := shutter{
		ioContext:     &io,
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}

	s.init()

	// Check up
	s.setCmd(1)
	s.tick()
	checkPins("req up ", t, io, HIGH, LOW)

	// Check wait
	s.setCmd(-1)
	for i := 0; i < s.DirSwitchWait; i++ {
		s.tick()
		checkPins("waiting ", t, io, LOW, LOW)
	}
	s.tick()

	// Check down
	checkPins("req down ", t, io, LOW, HIGH)

	s.tick()
	checkPins("req down ", t, io, LOW, LOW)
}
