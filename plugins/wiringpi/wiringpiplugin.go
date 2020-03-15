package main

import "gitlab.com/grill-tamasi/wscgo/plugins"

func GetAddons() []plugins.Addon {
	return []plugins.Addon{
		&mcp23017addon{},
		&pca9685addon{},
	}
}
