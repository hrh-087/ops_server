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
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(GetExecTimeMiddleware)
	// 游戏服安装
	mux.HandleFunc(task.InstallServerTypeName, task.HandleInstallServerTask)
	// 批量执行命令
	mux.HandleFunc(task.BatchCommandTypeName, task.HandleBatchCommand)
	// 更新游戏服镜像
	mux.HandleFunc(task.UpdateGameImageTypeName, task.HandleUpdateGameImage)
	// 关闭游戏服
	mux.HandleFunc(task.StopGameTypeName, task.HandleStopGame)
	// 更新配置文件
	mux.HandleFunc(task.UpdateGameJsonDataTypeName, task.HandleUpdateGameJsonData)
	// 启动游戏服
	mux.HandleFunc(task.StartGameTypeName, task.HandleStartGame)
	// 检查游戏服版本号
	mux.HandleFunc(task.CheckGameVersionTypeName, task.HandleCheckGameVersion)
	// 解压热更文件
	mux.HandleFunc(task.HotGameUnzipFileTypeName, task.HandleHotGameUnzipFile)
	// 同步热更文件到游戏服
	mux.HandleFunc(task.HotGameRsyncServerTypeName, task.HandleHotGameRsyncServer)
	// 同步热更文件到对应服务器
	mux.HandleFunc(task.HotGameRsyncHostTypeName, task.HandleHotGameRsyncHost)

	if err := srv.Run(mux); err != nil {
		global.OPS_LOG.Error("asynq task run failed", zap.Error(err))
	}
}
