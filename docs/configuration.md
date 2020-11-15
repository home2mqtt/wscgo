# Configuration

Configuration file uses ini format and is usually located at `/etc/wscgo.ini`. The path to the configuration file is given as argument to the 
invocation:

```sh
wscgo /etc/wscgo.ini
```

## MQTT

| Paramerer | Unit | Description |
| --- | --- | --- |
| host | url | URL of MQTT broker |
| user | string | (optional) username to use for authentication |
| password | string | (optional) password to use for authentication |
| clientid | string | (optional) client ID. Randomly generated if omitted |


```ini
[mqtt]
host = tcp://192.168.0.1:1883
user = username
password = password
clientid = uniqueclientid123
```

## IO extenders

### MCP23xxx

Supported members of MCP23xxx family: 
* I2C: MCP23008, MCP23009, MCP23016, MCP23017, MCP23018
* SPI: MCP23S08, MCP23S09, MCP23S17, MCP23S18

| Paramerer | Unit | Description |
| --- | --- | --- |
| address | int | (Only for I2C device) I2C address |

```ini
# Creates Pins MCP23017_<AddressHex>_PORT<A|B>_<0-7> e.g. MCP23017_20_PORTA_1
[mcp23017]
address = 0x20
```

```ini
# Creates Pins MCP23S17_PORT<A|B>_<0-7> e.g. MCP23S17_PORTA_1
[mcp23s17]
```

### PCA9685

```ini
# Initializes device with the given frequency and creates pins
# PCA9685_<AddressHex>_<0-15> e.g. PCA9685_40_0
[pca9685]
address = 0x40
frequency = 1000
```

## Supported devices

### Pin IDs

Pins are identified by string IDs as detected on the running platform or registered via IO extender configurations.

On raspberry Pi, IO headers available are detected.
* See [GPIO of BCM283x](https://godoc.org/periph.io/x/periph/host/bcm283x#Pin) 
* and [Headers exposed on RPi](https://godoc.org/periph.io/x/periph/host/rpi)

### Digital output

* [UC-1: Switch](usecases/uc-1.md)
* [UC-2: Switch (with IO extender)](usecases/uc-2.md)

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
pin = MCP23017_27_PORTA_0
```

Presented to Home assistant as a [switch](https://www.home-assistant.io/integrations/switch.mqtt)

### Digital input

* [UC-3: Button](usecases/uc-3.md)
* [UC-4: Button (with IO extender)](usecases/uc-4.md)

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to report state |
| pin | string | pin ID |

State payload
* ON
* OFF

Example:

```ini
[dinput:1]
name = Button
state_topic = home/f1/button1
pin = MCP23017_27_PORTA_0
```

Presented to Home assistant as a [binary sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt)

### Shutter

* [UC-5: Shutters](usecases/uc-5.md)
* [UC-6: Shutters (inverted output)](usecases/uc-6.md)


| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to receive commands |
| position_topic | string | Mqtt topic to report estimated position |
| uppin | string | pin ID to open shutter |
| downpin | string | pin ID to close shutter |
| dirswitchwait | int (100ms) | minimal time to wait between direction change |
| range | int (100ms) | Time needed to fully close/open shutter - used to estimate position |
| inverted | bool | (optional, default is false) indicates to use low-active logic for output |

Command payload:
* OPEN
* CLOSE
* OPENORSTOP - Stop if currently opening (since 0.5.4)
* CLOSEORSTOP - Stop if currently closing (since 0.5.4)
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
uppin = MCP23017_27_PORTA_0
downpin = MCP23017_27_PORTA_1
dirswitchwait = 20
range = 120
```

Presented to Home assistant as a [cover](https://www.home-assistant.io/integrations/cover.mqtt).

### Dimmable light

* [UC-7 Dimmed LED light with external PSU and galvanic isolation](usecases/uc-7.md) 
* [UC-8 Dimmed LED light with common ground](usecases/uc-8.md)

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| topic | string | Mqtt topic to receive commands |
| onpin | string | (optional) pin ID to enable or disable power source of the light |
| pwmpin | string | pin ID to use as PWM output |
| ondelay | int (100 ms) | (optional, default is 0) Delay after setting onpin to high and before raising pwm output to allow the external PSU to settle |
| inverted | bool | (optional, default is false) if set, PWM output configured as low-active |
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

### Comet Blue

* [UC-10 Blue Comet thermostatic valve](usecases/uc-10.md)

| Paramerer | Unit | Description |
| --- | --- | --- |
| name | string | Device name |
| mac | string | BLE MAC address of the thermostat |
| duration | int (100ms) | Time to wait between reading state of thermostat |
| topic | string | Base MQTT topic for this device |

```ini
[cometblue:1]
name = Comet Blue 1
mac = 11:22:33:44:55:66
duration = 600
topic = home/cb1
```

Presented to Home assistant as a [HVAC device](https://www.home-assistant.io/integrations/climate.mqtt/).

## Using with Home Assistant

All configured devices publish proper discovery configuration upon startup with `homeassistant/` prefix. Additionally, sending discovery
messages can be triggered via topic `discovery` from homeassistant:

![discover](discover.png)
