# Building wscgo

Supported Go version is >=1.13

```bash
export GOOS=linux
export GOARCH=arm
# raspberry pi zero is compatible with ARMv5
export GOARM=5
# version string is set via linker flag
go build -v -ldflags "-X main.Version=0.5.2-snapshot" -trimpath -o wscgo ./cmd/wscgo
```