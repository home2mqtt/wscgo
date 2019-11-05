package main

type ioContext interface {
	digitalWrite(pin int, value int)
	digitalRead(pin int) int
	pinMode(pin int, mode int)
}

type io struct {
	ioContext
	out   int
	value int
	topic string
}

func (io *io) tick() {
	io.digitalWrite(io.out, io.value)
}

func (io *io) init() {
	io.pinMode(io.out, OUTPUT)
	io.value = LOW
}
