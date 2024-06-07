package main

import (
	"log"
	"os"

	"github.com/home2mqtt/wscgo/config"
	"periph.io/x/host/v3"

	"github.com/home2mqtt/wscgo/integration"
)

func main() {
	log.Println("wscgo version ", Version)
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("usage: wscgo config.ini")
	}
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}
	conf := config.LoadConfig(args[0])

	instance := &integration.WscgoInstance{
		Version: Version,
	}
	instance.Configure(conf)
	instance.Start()
	instance.Loop()
}
