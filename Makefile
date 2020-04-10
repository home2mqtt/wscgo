all: wscgo

install: wscgo
	cp wscgo /usr/bin/
	cp wscgo.ini /etc/
	cp wscgo.service /etc/systemd/system/

wscgo:
	go build -v -ldflags "-X main.Version=${VERSION}" -trimpath -o wscgo ./cmd/wscgo