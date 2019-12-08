# wscgo

wscgo is a highly configurable controller intended for home automation written in Go. Supported 
platforms include [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) 
and [Orange Pi Zero](http://www.orangepi.org/orangepizero/), it should work on any device where 
[Wiring PI](http://wiringpi.com/) can be compiled.

### Main features

* Uses [Wiring PI](http://wiringpi.com/)
* Capable using I/O extenders (e.g. [MCP23017](https://www.microchip.com/wwwproducts/en/MCP23017))
* [MQTT](http://mqtt.org/)-based protocol compatible with [Home Assistant](https://www.home-assistant.io/integrations/mqtt/) 
* Supported devices:
  * Digital output (switch)
  * Digital input
  * Dimmable light
  * Window shutters
