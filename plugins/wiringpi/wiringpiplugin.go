package main

import "github.com/balazsgrill/wscgo/plugins"

func GetAddons() []plugins.Addon {
	return []plugins.Addon{
		&mcp23017addon{},
		&pca9685addon{},
	}
}
