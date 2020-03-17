# wscgo

wscgo is a highly configurable controller intended for home automation written in Go. Supported 
platforms include [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) 
and [Orange Pi Zero](http://www.orangepi.org/orangepizero/), it should work on any device supported by [periph.io](https://periph.io/

![Place of wscgo](https://raw.githubusercontent.com/wiki/balazsgrill/wscgo/place-of-wscgo.png)

See [user guide](https://github.com/balazsgrill/wscgo/wiki/User-guide)

### Main features

* Uses [periph.io](https://periph.io/)
  * Also can use [Wiring PI](http://wiringpi.com/)
* Capable using I/O extenders (e.g. [MCP23017](https://www.microchip.com/wwwproducts/en/MCP23017))
* [MQTT](http://mqtt.org/)-based protocol compatible with [Home Assistant](https://www.home-assistant.io/integrations/mqtt/) 
* Supported devices:
  * Digital output (switch)
  * Digital input
  * Dimmable light
  * Window shutters
