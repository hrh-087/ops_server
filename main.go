package main

import (
	"go.uber.org/zap"
	"ops-server/core"
	"ops-server/global"
	"ops-server/initialize"
	"ops-server/job/inspector"
	"ops-server/job/scheduler"
	"ops-server/job/server"
	"ops-server/job/workers"
)

func main() {
	// 初始化Viper 读取配置文件
	global.OPS_VP = core.Viper()
	// 初始化jwt过期时间
	//initialize.OtherInit()
	// 初始化zap日志库
	global.OPS_LOG = core.Zap()
	zap.ReplaceGlobals(global.OPS_LOG)
	// gorm连接数据库
	global.OPS_DB = initialize.Gorm()
	// 异步队列
	global.AsynqInspect = inspector.NewAsynqInspector()
	global.AsynqClient = server.NewAsynqClinet()
	global.AsynqScheduler = scheduler.NewAsynqScheduler()

	if global.OPS_DB != nil {
		// 初始化表
		initialize.RegisterTables()
		// 程序结束前关闭数据库链接
		db, _ := global.OPS_DB.DB()
		defer db.Close()
	}

	// 运行asynq消费者
	go func() {
		workers.InitWorkers()
	}()

	core.RunWindowsServer()
}
