# UC-2 Driving load via relay connected to IO extender

![Direct relay connection](relay_ioex_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

[plugin]
path = /usr/local/lib/wscgo-wpi-rpizw.so

[mcp23017]
address = 0x20
expansionBase = 100

[switch:0]
name = Relay
topic = home/relay/0
pin = 100
```

> Relay by [arduinomodules](https://arduinomodules.info/ky-019-5v-relay-module/)\
> Bulb by [alfreddagenais](https://github.com/alfreddagenais/fritzing-components/)