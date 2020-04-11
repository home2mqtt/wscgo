all: wscgo

install: wscgo
	sudo cp wscgo /usr/bin/
	sudo cp wscgo.ini /etc/
	sudo cp wscgo.service /etc/systemd/system/

wscgo:
	go build -v -ldflags "-X main.Version=${VERSION}" -trimpath -o wscgo ./cmd/wscgo