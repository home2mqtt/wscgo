package main

type nand struct {
	ioContext
	in1 int
	in2 int
	out int
}

func (nand *nand) tick() {
	in1 := nand.digitalRead(nand.in1)
	in2 := nand.digitalRead(nand.in2)
	out := LOW
	if !((in1 == HIGH) && (in2 == HIGH)) {
		out = HIGH
	}
	nand.digitalWrite(nand.out, out)
}

func (nand *nand) init() {
	nand.pinMode(nand.in1, INPUT)
	nand.pinMode(nand.in2, INPUT)
	nand.pinMode(nand.out, OUTPUT)
}
