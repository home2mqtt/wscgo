package main

import (
	"bufio"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	goio "io"
	"log"
	"strconv"
)

type serialconf struct {
	baudrate  uint
	portname  string
	topicroot string

	light         int
	port          goio.ReadWriteCloser
	valuecallback func(int, int)
}

func (serialconf *serialconf) setLight(value int) {
	serialconf.light = value
}

func (serialconf *serialconf) connect() {
	options := serial.OpenOptions{
		PortName:        serialconf.portname,
		BaudRate:        serialconf.baudrate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	serialconf.port = port
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(port)
		for scanner.Scan() {
			var arg int
			var value int
			count, err := fmt.Sscanf(scanner.Text(), "!%d %d", &arg, &value)
			if (err == nil) && (count == 2) {
				if serialconf.valuecallback != nil {
					serialconf.valuecallback(arg, value)
				}
			}
		}
	}()
}

func (serialconf *serialconf) init() {
	serialconf.connect()
}

func (serialconf *serialconf) tick() {
	_, err := serialconf.port.Write([]byte(strconv.Itoa(serialconf.light)))
	if err != nil {
		log.Print("Serial send error: ", err)
		serialconf.port.Close()
		serialconf.connect()
	}
}
