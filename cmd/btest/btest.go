package main

import (
	"encoding/hex"
	"log"

	"github.com/balazsgrill/wscgo/devices/cometblue"
)

func main() {
	cb, err := cometblue.Dial("E9:33:A1:84:05:6A")
	if err != nil {
		log.Fatal(err)
	}
	defer cb.Close()
	err = cb.Authenticate()
	if err != nil {
		log.Fatal(err)
	}
	bstr, err := cb.ReadBatteryRaw()
	if err != nil {
		log.Fatal(err)
	}
	str := hex.EncodeToString(bstr)
	log.Println(str)
	bstr, err = cb.ReadTemperaturesRaw()
	if err != nil {
		log.Fatal(err)
	}
	str = hex.EncodeToString(bstr)
	log.Println(str)
}
