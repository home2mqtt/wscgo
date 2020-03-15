package main

import "gitlab.com/grill-tamasi/wscgo/plugins"

func GetAddons() []plugins.Addon {
	return {
		&mcp23017addon{},
		&pcapca9685addon{},
	}
}
