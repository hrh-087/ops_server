package core

import (
	"fmt"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/initialize"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.OPS_CONFIG.System.UseRedis || global.OPS_CONFIG.System.UseMultipoint {
		//initialize.Gorm()
		initialize.Redis()
	}

	if global.OPS_CONFIG.System.UseMongo {
		//err := initialize.Mongo.Initialization()
		//if err != nil {
		//	zap.L().Error(fmt.Sprintf("%+v", err))
		//}
	}

	if global.OPS_DB != nil {
	}

	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.OPS_CONFIG.System.Addr)

	s := initServer(address, Router)

	global.OPS_LOG.Info("server run success on ", zap.String("address", address))

	global.OPS_LOG.Error(s.ListenAndServe().Error())
}
