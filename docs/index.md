# WSCGO

DIY home-automation using cheap off-the-shelf hardware components.

> :warning: Disclaimer: Some applications involve working with high voltage (230VAC). Consult with an electrician if you haven't got the expertise, the required safety measures are exluded from this documentation.

- [Installation](install.md)
- [Use cases](usecases/usecaselist.md)
- [Configuration reference](configuration.md)
- [Build from source](build.md)

## What?

So you have some appliance (e.g. a motorized [shutters](usecases/uc-5.md), [lights](usecases/uc-1.md) or some [buttons](usescases/uc-3.md)) in your home which you want to automate using home automation software (e.g. [homeassistant](https://www.home-assistant.io/)). You may have an SBC ([raspberry pi](https://www.raspberrypi.org/), [orange pi](http://www.orangepi.org/)) in your drawer along with some peripheral boards (Relay boards, IO extenders). You can connect these components with some dupont wire but you still need some software on your SBC to make use of them. 

Instead of writing a software from scratch specialized for your particular setup, wscgo can be [configured](configuration.md) for a wide range of [uses](usecases/usecaselist.md).

## Why?

The development is mainly driven by [my](https://github.com/balazsgrill) needs regarding home automation.

### Why SBCs?

Why not some cheaper alternative, like [ESP](https://en.wikipedia.org/wiki/ESP8266)? Currently I have five controllers in my house. I could save a few bucks by using cheaper devices, but having linux on these computers I gain ssh access, remote monitoring and easy [updates via package manager](install.md).

## How?

Wscgo written in [Go](https://golang.org) and uses [Periph](https://periph.io) for low-level device access. [WiringPi](http://wiringpi.com/) is also used for now, but it will be removed in the future.

## License

Wscgo is open source, released under [GPL3](../LICENSE). Contributions are welcome.