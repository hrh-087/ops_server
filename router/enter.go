package router

import "ops-server/router/system"

type routerGroup struct {
	System system.RouterGroup
}

var RouterGroupApp = new(routerGroup)
