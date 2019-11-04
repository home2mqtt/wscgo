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
