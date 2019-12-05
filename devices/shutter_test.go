package devices

import (
	"testing"

	"gitlab.com/grill-tamasi/wscgo/tests"
	"gitlab.com/grill-tamasi/wscgo/wiringpi"
)

func checkPins(msg string, t *testing.T, io *tests.TestIo, up bool, down bool) {
	if io.Values[0] != up || io.Values[1] != down {
		t.Errorf("%s UP[exp-actal]: %t - %t, DOWN[exp-actal]: %t - %t\n", msg, up, io.Values[0], down, io.Values[1])
	}
}

func createShutterForTest() (*shutter, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	sc := ShutterConfig{
		UpPin:         0,
		DownPin:       1,
		DirSwitchWait: 20,
		Range:         10,
	}
	s, _ := CreateShutter(io, &sc).(*shutter)
	return s, io
}

func TestInit(t *testing.T) {
	s, io := createShutterForTest()

	s.Initialize()
	if io.Modes[0] != wiringpi.OUTPUT {
		t.Fatal()
	}
	if io.Modes[1] != wiringpi.OUTPUT {
		t.Fatal()
	}
}

func TestUp(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	s.Open()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg up", t, io, wiringpi.HIGH, wiringpi.LOW)
	}
	s.Tick()
	checkPins("reg stop", t, io, wiringpi.LOW, wiringpi.LOW)
}

func TestDown(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	s.Close()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg down", t, io, wiringpi.LOW, wiringpi.HIGH)
	}
	s.Tick()
	checkPins("reg stop", t, io, wiringpi.LOW, wiringpi.LOW)
}

func TestStop(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	for i := 0; i < 50; i++ {
		s.Tick()
		checkPins("req zero ", t, io, wiringpi.LOW, wiringpi.LOW)
	}

}

func TestDirectionChange(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, wiringpi.HIGH, wiringpi.LOW)

	// Check wait
	s.setCmd(-1)
	for i := 0; i < s.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, wiringpi.LOW, wiringpi.LOW)
	}
	s.Tick()

	// Check down
	checkPins("req down ", t, io, wiringpi.LOW, wiringpi.HIGH)

	s.Tick()
	checkPins("req down ", t, io, wiringpi.LOW, wiringpi.LOW)
}

func TestDirectionChangeWithExtraWait(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, wiringpi.HIGH, wiringpi.LOW)

	// Extra wait
	for i := 0; i < s.DirSwitchWait*2; i++ {
		s.Tick()
		checkPins("waiting ", t, io, wiringpi.LOW, wiringpi.LOW)
	}
	s.Tick()
	s.Tick()

	// Check down
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, wiringpi.LOW, wiringpi.HIGH)

	s.Tick()
	checkPins("req stop ", t, io, wiringpi.LOW, wiringpi.LOW)
}

func TestDirectionChangeWithStop(t *testing.T) {
	s, io := createShutterForTest()
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, wiringpi.HIGH, wiringpi.LOW)

	// Stop --> wait
	s.setCmd(0)
	for i := 0; i < s.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, wiringpi.LOW, wiringpi.LOW)
	}

	// Check wait
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, wiringpi.LOW, wiringpi.HIGH)

	s.Tick()
	checkPins("req down ", t, io, wiringpi.LOW, wiringpi.LOW)
}
