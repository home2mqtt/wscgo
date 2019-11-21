package devices

import "testing"

const HIGH bool = true
const LOW bool = false

type testIo struct {
	modes  []bool
	values []bool
}

func (io *testIo) DigitalWrite(pin int, value bool) {
	io.values[pin] = value
}

func (io *testIo) DigitalRead(pin int) bool {
	return io.values[pin]
}

func (io *testIo) PinMode(pin int, mode bool) {
	io.modes[pin] = mode
}

func checkPins(msg string, t *testing.T, io *testIo, up bool, down bool) {
	if io.values[0] != up || io.values[1] != down {
		t.Errorf("%s UP[exp-actal]: %t - %t, DOWN[exp-actal]: %t - %t\n", msg, up, io.values[0], down, io.values[1])
	}
}

func createShutterForTest() (*shutter, *testIo) {
	io := testIo{
		modes:  []bool{false, false, false},
		values: []bool{false, false, false},
	}
	sc := ShutterConfig{
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}
	s, _ := CreateShutter(&io, &sc).(*shutter)
	return s, &io
}

func TestInit(t *testing.T) {
	s, io := createShutterForTest()

	s.Initialize()
	if !io.modes[0] {
		t.Fatal()
	}
	if !io.modes[1] {
		t.Fatal()
	}
}

func TestUp(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	s.setCmd(10)
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg up", t, io, HIGH, LOW)
	}
	s.Tick()
	checkPins("reg stop", t, io, LOW, LOW)
}

func TestDown(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	s.setCmd(-10)
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg down", t, io, LOW, HIGH)
	}
	s.Tick()
	checkPins("reg stop", t, io, LOW, LOW)
}

func TestStop(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	for i := 0; i < 50; i++ {
		s.Tick()
		checkPins("req zero ", t, io, LOW, LOW)
	}

}

func TestDirectionChange(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, HIGH, LOW)

	// Check wait
	s.setCmd(-1)
	for i := 0; i < s.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, LOW, LOW)
	}
	s.Tick()

	// Check down
	checkPins("req down ", t, io, LOW, HIGH)

	s.Tick()
	checkPins("req down ", t, io, LOW, LOW)
}

func TestDirectionChangeWithExtraWait(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, HIGH, LOW)

	// Extra wait
	for i := 0; i < s.DirSwitchWait*2; i++ {
		s.Tick()
		checkPins("waiting ", t, io, LOW, LOW)
	}
	s.Tick()
	s.Tick()

	// Check down
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, LOW, HIGH)

	s.Tick()
	checkPins("req stop ", t, io, LOW, LOW)
}

func TestDirectionChangeWithStop(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, HIGH, LOW)

	// Stop --> wait
	s.setCmd(0)
	for i := 0; i < s.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, LOW, LOW)
	}

	// Check wait
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, LOW, HIGH)

	s.Tick()
	checkPins("req down ", t, io, LOW, LOW)
}
