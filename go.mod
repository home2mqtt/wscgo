module gitlab.com/grill-tamasi/wscgo

go 1.13

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	golang.org/x/net v0.0.0-20191101175033-0deb6923b6d9 // indirect
	gopkg.in/ini.v1 v1.49.0
	periph.io/x/periph v3.6.2+incompatible
	gitlab.com/grill-tamasi/wscgo/plugins v0.5.0
)

replace gitlab.com/grill-tamasi/wscgo/plugins => ./plugins
