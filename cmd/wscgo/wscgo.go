package main

import (
	"log"
	"os"

	"github.com/balazsgrill/wscgo/config"
	"github.com/balazsgrill/wscgo/integration"
	"periph.io/x/periph/host"

	_ "github.com/balazsgrill/wscgo/integration"
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
