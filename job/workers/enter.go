package workers

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/job/task"
)

func InitWorkers() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     global.OPS_CONFIG.Asynq.Addr + ":" + global.OPS_CONFIG.Asynq.Port,
			Password: global.OPS_CONFIG.Asynq.Pass,
			DB:       global.OPS_CONFIG.Asynq.Db,
		},
		asynq.Config{
			Concurrency: global.OPS_CONFIG.Asynq.Concurrency,
			//Logger: global.OPS_LOG,
			Queues: map[string]int{
				"default": 1,
				"cron":    2,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(GetExecTimeMiddleware)

	gameMux := asynq.NewServeMux()
	gameMux.Use(GetExecTimeMiddleware)

	cronMux := asynq.NewServeMux()
	cronMux.Use(CronMiddleware)

	// 公共
	// 初始化项目
	mux.HandleFunc(task.InitProjectTypeName, task.HandleInitProject)

	// 游戏服
	// 游戏服安装
	gameMux.HandleFunc(task.InstallGameServerTypeName, task.HandleInstallServer)
	// 批量执行命令
	gameMux.HandleFunc(task.BatchCommandTypeName, task.HandleBatchCommand)
	// 更新游戏服镜像
	gameMux.HandleFunc(task.UpdateGameImageTypeName, task.HandleUpdateGameImage)
	// 关闭游戏服
	gameMux.HandleFunc(task.StopGameTypeName, task.HandleStopGame)
	// 更新配置文件
	gameMux.HandleFunc(task.RsyncGameJsonConfigTypeName, task.HandleRsyncGameJsonConfig)
	// 启动游戏服
	gameMux.HandleFunc(task.StartGameTypeName, task.HandleStartGame)
	// 检查游戏服版本号
	gameMux.HandleFunc(task.CheckGameVersionTypeName, task.HandleCheckGameImageVersion)
	// 解压热更文件
	gameMux.HandleFunc(task.HotGameUnzipFileTypeName, task.HandleHotGameUnzipFile)
	// 同步热更文件到游戏服
	gameMux.HandleFunc(task.HotGameRsyncServerTypeName, task.HandleHotGameRsyncServer)
	// 同步热更文件到对应服务器
	gameMux.HandleFunc(task.HotGameRsyncHostTypeName, task.HandleHotGameRsyncHost)
	// 同步游戏服配置文件
	gameMux.HandleFunc(task.RsyncGameConfigTypeName, task.HandleRsyncGameConfig)
	// 同步游戏服脚本
	gameMux.HandleFunc(task.RsyncGameScriptTypeName, task.HandleRsyncGameScript)

	// cron
	cronMux.HandleFunc(task.CronCloseMatchBlockTypeName, task.HandleCloseMatchBlock) // 关闭匹配服
	cronMux.HandleFunc(task.CronKickPlayerTypeName, task.HandleKickPlayer)

	mux.Handle("cron:", cronMux) // 匹配所有cron:开头的task
	mux.Handle("game:", gameMux) // 匹配所有game:开头的task

	if err := srv.Run(mux); err != nil {
		global.OPS_LOG.Error("asynq task run failed", zap.Error(err))
	}

}
