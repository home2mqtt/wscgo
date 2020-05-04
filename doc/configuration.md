# Configuration

Configuration file uses ini format and is usually located at `/etc/wscgo.ini`. The path to the configuration file is given as argument to the 
invocation:

```sh
wscgo /etc/wscgo.ini
```

## MQTT

```ini
[mqtt]
host = tcp://192.168.0.1:1883
user = username
password = password
clientid = uniqueclientid123
```

## IO extenders

### WiringPi plugin

IO extenders are accessible via platform-dependent plugin using WiringPi.

```ini
# Rasperry Pi Zero
[plugin]
path = /usr/local/lib/wscgo-wpi-rpizw.so

# Orange Pi Zero
[plugin]
path = /usr/local/lib/wscgo-wpi-opiz.so
```

### MCP23017

```ini
# Creates Pins 100-115
[mcp23017]
address = 0x20
expansionBase = 100
```

## Supported devices

### Digital output

* [UC-1: Switch](usecases/uc-1/uc-1.md)
* [UC-2: Switch (with IO extender)](usecases/uc-2/uc-2.md)

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to receive commands |
| pin | int | wiringpi pin |

Command payload
* ON
* OFF

Example: 

```ini
[switch:0]
name = Entry hall light
topic = home/f1/entry/light
pin = 101
```

Presented to Home assistant as a [switch](https://www.home-assistant.io/integrations/switch.mqtt)

### Digital input

* [UC-3: Button](usecases/uc-3/uc-3.md)
* [UC-4: Button (with IO extender)](usecases/uc-4/uc-4.md)

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to report state |
| pin | int | wiringpi pin |

State payload
* ON
* OFF

Example:

```ini
[dinput:1]
name = Button
state_topic = home/f1/button1
pin = 102
```

Presented to Home assistant as a [binary sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt)

### Shutter

* [UC-5: Shutters](usecases/uc-5/uc-5.md)
* [UC-6: Shutters (inverted output)](usecases/uc-6/uc-6.md)


| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to receive commands |
| position_topic | string | Mqtt topic to report estimated position |
| uppin | int | pin to open shutter |
| downpin | int | pin to close shutter |
| dirswitchwait | int (100ms) | minimal time to wait between direction change |
| range | int (100ms) | Time needed to fully close/open shutter - used to estimate position |
| opt_groupTopic | string | (optional) additional topic intended to receive commands. The same topic may be set to multple shutters |

Command payload:
* OPEN
* CLOSE
* STOP
* <integer>
  * Negative - move down by given time (in 100 ms)
  * Positive - move up by given time (in 100ms)
  * Zero - stop

Position payload: estimated position between 0 (fully closed) and range (fully open)

Example:

```ini
[shutter:4]
name = Window 1
topic = home/f1/shutter/4
position_topic = home/f1/shutter/4/state
uppin = 103
downpin = 104
dirswitchwait = 20
range = 120
```

Presented to Home assistant as a [cover](https://www.home-assistant.io/integrations/cover.mqtt).

### Dimmable light

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to receive commands |
| onpin | int | (optional) wiringpi pin to enable or disable light |
| pwmpin | int | wiringpi pin to use as PWM output |
| ondelay | int (100 ms) | (optional, default is 0) Delay after setting onpin to high and before raising pwm output to allow the external PSU to settle |
| inverted | bool | (optional, default is false) if set, PWM output is inverted |
| speed | int | Maximum PWM value change by tick (100ms) |

Command payload:
* ON: Set maximum brightness
* OFF: Set brightness to 0 then turn light off
* <integer>: 0-1023 set brightness

```ini
[light:6]
name = Light
topic = home/f2/light
pwmpin = 7
speed = 20
ondelay = 1
inverted=true
```

Presented to Home assistant as a [light](https://www.home-assistant.io/integrations/light.mqtt).

#### Note on PWM output

Wscgo tries to determine whether hardware PWM is available on the used pin, otherwise it falls back to software PWM. Via IO extenders, software PWM is disabled.

## Using with Home Assistant

All configured devices publish proper discovery configuration upon startup with `homeassistant/` prefix. Additionally, sending discovery
messages can be triggered via topic `discovery` from homeassistant:

![discover](discover.png)
