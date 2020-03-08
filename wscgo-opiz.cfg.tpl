Section: misc
Priority: optional
# Homepage: <enter URL here; no default>
Standards-Version: 3.9.2

Package: wscgo-opiz
Version: ${VERSION}
Maintainer: Balázs Grill <balazs.grill@live.com>
# Pre-Depends: <comma-separated list of packages>
Depends:
# Recommends: <comma-separated list of packages>
# Suggests: <comma-separated list of packages>
# Provides: <comma-separated list of packages>
# Replaces: <comma-separated list of packages>
Postinst: postinst
Postrm: postrm
Architecture: armhf
# Multi-Arch: <one of: foreign|same|allowed>
# Copyright: <copyright file; defaults to GPL2>
# Changelog: <changelog file; defaults to a generic changelog>
# Readme: <README.Debian file; defaults to a generic one>
# Extra-Files: <comma-separated list of additional files for the doc directory>
Files: wscgo /usr/bin/
 wscgo.ini /etc/
 wscgo.service /etc/systemd/system/
 libwiringPi.so.2.0 /usr/local/lib/
 libwiringPi.so /usr/local/lib/
#  <more pairs, if there's more than one file to include. Notice the starting space>
Description: Window-shutter controller
 Configurable home automation controller written in Go