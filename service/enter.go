package service

import "ops-server/service/system"

type serviceGroup struct {
	SystemServiceGroup system.ServiceGroup
}

var ServiceGroupApp = new(serviceGroup)
