# UC-4 Digial input via IO extender

![Direct relay connection](button_ioex_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

[mcp23017]
address = 0x20

[dinput:0]
name = Button
state_topic = home/relay/0
pin = MCP23017_20_PORTA_0
```
