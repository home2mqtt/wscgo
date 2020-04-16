package devices

import (
	"testing"

	"github.com/balazsgrill/wscgo/tests"
	"periph.io/x/periph/conn/gpio"
)

func checkPins(msg string, t *testing.T, io *tests.TestIo, up gpio.Level, down gpio.Level) {
	if io.Pins[0].L != up || io.Pins[1].L != down {
		t.Errorf("%s UP[exp-actal]: %t - %t, DOWN[exp-actal]: %t - %t\n", msg, up, io.Pins[0].L, down, io.Pins[1].L)
	}
}

func createShutterForTest(inverted bool) (*shutter, *tests.TestIo) {
	io := tests.CreateTestIo(3)
	sc := ShutterConfig{
		UpPin:         "Test_0",
		DownPin:       "Test_1",
		DirSwitchWait: 20,
		Range:         10,
		Inverted:      inverted,
	}
	is, _ := CreateShutter(&sc)
	s, _ := is.(*shutter)
	return s, io
}

func TestInit(t *testing.T) {
	s, io := createShutterForTest(false)

	s.Initialize()
	if io.Pins[0].L != gpio.Low {
		t.Fatal()
	}
	if io.Pins[1].L != gpio.Low {
		t.Fatal()
	}
}

func TestInitInverted(t *testing.T) {
	s, io := createShutterForTest(true)

	s.Initialize()
	if io.Pins[0].L != gpio.High {
		t.Fatal()
	}
	if io.Pins[1].L != gpio.High {
		t.Fatal()
	}
}

func TestUp(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	s.Open()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg up", t, io, gpio.High, gpio.Low)
	}
	s.Tick()
	checkPins("reg stop", t, io, gpio.Low, gpio.Low)
}

func TestUpInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	s.Open()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg up", t, io, gpio.Low, gpio.High)
	}
	s.Tick()
	checkPins("reg stop", t, io, gpio.High, gpio.High)
}

func TestDown(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	s.Close()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg down", t, io, gpio.Low, gpio.High)
	}
	s.Tick()
	checkPins("reg stop", t, io, gpio.Low, gpio.Low)
}

func TestDownInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	s.Close()
	for i := 0; i < 10; i++ {
		s.Tick()
		checkPins("reg down", t, io, gpio.High, gpio.Low)
	}
	s.Tick()
	checkPins("reg stop", t, io, gpio.High, gpio.High)
}

func TestStop(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	for i := 0; i < 50; i++ {
		s.Tick()
		checkPins("req zero ", t, io, gpio.Low, gpio.Low)
	}

}

func TestStopInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	for i := 0; i < 50; i++ {
		s.Tick()
		checkPins("req zero ", t, io, gpio.High, gpio.High)
	}

}

func TestDirectionChange(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.High, gpio.Low)

	// Check wait
	s.setCmd(-1)
	for i := 0; i < s.config.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.Low, gpio.Low)
	}
	s.Tick()

	// Check down
	checkPins("req down ", t, io, gpio.Low, gpio.High)

	s.Tick()
	checkPins("req down ", t, io, gpio.Low, gpio.Low)
}

func TestDirectionChangeInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.Low, gpio.High)

	// Check wait
	s.setCmd(-1)
	for i := 0; i < s.config.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.High, gpio.High)
	}
	s.Tick()

	// Check down
	checkPins("req down ", t, io, gpio.High, gpio.Low)

	s.Tick()
	checkPins("req down ", t, io, gpio.High, gpio.High)
}

func TestDirectionChangeWithExtraWait(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.High, gpio.Low)

	// Extra wait
	for i := 0; i < s.config.DirSwitchWait*2; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.Low, gpio.Low)
	}
	s.Tick()
	s.Tick()

	// Check down
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, gpio.Low, gpio.High)

	s.Tick()
	checkPins("req stop ", t, io, gpio.Low, gpio.Low)
}

func TestDirectionChangeWithExtraWaitInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.Low, gpio.High)

	// Extra wait
	for i := 0; i < s.config.DirSwitchWait*2; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.High, gpio.High)
	}
	s.Tick()
	s.Tick()

	// Check down
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, gpio.High, gpio.Low)

	s.Tick()
	checkPins("req stop ", t, io, gpio.High, gpio.High)
}

func TestDirectionChangeWithStop(t *testing.T) {
	s, io := createShutterForTest(false)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.High, gpio.Low)

	// Stop --> wait
	s.setCmd(0)
	for i := 0; i < s.config.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.Low, gpio.Low)
	}

	// Check wait
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, gpio.Low, gpio.High)

	s.Tick()
	checkPins("req down ", t, io, gpio.Low, gpio.Low)
}

func TestDirectionChangeWithStopInverted(t *testing.T) {
	s, io := createShutterForTest(true)
	s.Initialize()

	// Check up
	s.setCmd(1)
	s.Tick()
	checkPins("req up ", t, io, gpio.Low, gpio.High)

	// Stop --> wait
	s.setCmd(0)
	for i := 0; i < s.config.DirSwitchWait; i++ {
		s.Tick()
		checkPins("waiting ", t, io, gpio.High, gpio.High)
	}

	// Check wait
	s.setCmd(-1)
	s.Tick()
	checkPins("req down ", t, io, gpio.High, gpio.Low)

	s.Tick()
	checkPins("req down ", t, io, gpio.High, gpio.High)
}
