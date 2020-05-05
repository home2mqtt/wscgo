# UC-6 Shutter connection via IO extender and relay board

Some relay boards<sup>[(1)][1]</sup> have low-active inputs which require inverted output logic.

> :warning: Low-active modules requires driven sink output, don't connect 5V boards directly to Raspberry Pi!

![Shutter connection schematic](./shutter_connection_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

[plugin]
path = /usr/local/lib/wscgo-wpi-rpizw.so

[mcp23017]
address = 0x20
expansionBase = 100

[shutter:0]
name = testShutter
topic = test/shutter/0
position_topic = test/shutter/0/state
uppin = 100
downpin = 101
dirswitchwait = 20
range = 120
inverted = true
```

[1]: https://arduinodiy.wordpress.com/2018/09/04/the-16-relay-module-and-the-raspberry-pi-not-an-ideal-marriage/