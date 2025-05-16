package main

import (
	"flag"
	"go.uber.org/zap"
	"ops-server/core"
	"ops-server/global"
	"ops-server/initialize"
	"ops-server/job/inspector"
	"ops-server/job/scheduler"
	"ops-server/job/server"
	"ops-server/job/workers"
	_ "ops-server/source/system" // 加载初始数据
)

func main() {
	var config string
	var operate string

	flag.StringVar(&config, "c", "", "选择配置文件")
	flag.StringVar(&operate, "o", "", "执行操作 worker 开启异步任务 scheduler 开启定时任务 initData 加载初始化数据 exportData 导出数据")
	flag.Parse()

	// 初始化Viper 读取配置文件
	global.OPS_VP = core.Viper(config)
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

	switch operate {
	case "worker":
		// 运行asynq消费者
		workers.InitWorkers()
	case "scheduler":
		// 运行定时任务消费者
		scheduler.InitScheduler()
	case "initData":
		err := initialize.InitDBServiceApp.InitData()
		if err != nil {
			global.OPS_LOG.Error("初始化数据失败", zap.Error(err))
			panic(err)
		}
	case "exportData":
		// 导出api跟菜单数据
		err := initialize.ExportData()
		if err != nil {
			global.OPS_LOG.Error("导出数据失败", zap.Error(err))
			panic(err)
		}
	default:
		core.RunWindowsServer()
	}

}
