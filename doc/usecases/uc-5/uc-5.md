# UC-5 Shutter connection via IO extender and relay board

![Shutter connection schematic](./shutter_connection_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

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
```