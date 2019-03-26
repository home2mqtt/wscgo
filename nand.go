package main

type nand struct {
	in1 int
	in2 int
	out int
}

func (nand *nand) tick() {
	in1 := digitalRead(nand.in1)
	in2 := digitalRead(nand.in2)
	out := LOW
	if !((in1 == HIGH) && (in2 == HIGH)) {
		out = HIGH
	}
	digitalWrite(nand.out, out)
}

func (nand *nand) init() {
	pinMode(nand.in1, INPUT)
	pinMode(nand.in2, INPUT)
	pinMode(nand.out, OUTPUT)
}
