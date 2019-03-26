package main

type io struct {
	out   int
	value int
	topic string
}

func (io *io) tick() {
	digitalWrite(io.out, io.value)
}

func (io *io) init() {
	pinMode(io.out, OUTPUT)
	io.value = LOW
}
