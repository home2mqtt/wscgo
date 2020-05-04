# Install from debian repository

Applicable for debian-based systems (raspbian, armbian)

```
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
echo "deb https://dl.bintray.com/balazsgrill/wscgo unstable main" | sudo tee -a /etc/apt/sources.list
sudo apt update
```

Install
```
sudo apt install wscgo
```

Install wiringpi binary and plugin for appropriate platform:
* Orange Pi Zero: `sudo apt install wscgo-opiz`
* Raspberry Pi Zero W: `sudo apt install wscgo-rpizw`
