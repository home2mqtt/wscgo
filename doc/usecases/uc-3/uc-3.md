# UC-3 Digial input

![Direct relay connection](button_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

[dinput:0]
name = Button
state_topic = home/relay/0
pin = GPIO23
```