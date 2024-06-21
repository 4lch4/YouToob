package templates

import (
	_ "embed"
)

type BaseRouteHelp struct {
	BaseURL     string
	ChannelName string
}

//go:embed Base-Route-Help.tmpl
var baseRouteHelp string

func GetBaseRouteHelp() string {
	return baseRouteHelp
}
