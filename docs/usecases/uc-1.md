# UC-1 Driving load via relay

![Direct relay connection](relay_bb.png)

## Configuration

```ini
[mqtt]
host = tcp://127.0.0.1:1883

[switch:0]
name = Relay
topic = home/relay/0
pin = GPIO23
```