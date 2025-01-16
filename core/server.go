package core

import (
	"fmt"
	"log"
	"ops-server/global"
	"ops-server/initialize"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.OPS_CONFIG.System.UseRedis || global.OPS_CONFIG.System.UseMultipoint {
		//initialize.Gorm()
	}

	if global.OPS_DB != nil {
	}

	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.OPS_CONFIG.System.Addr)

	s := initServer(address, Router)

	log.Fatal(s.ListenAndServe().Error())
}
