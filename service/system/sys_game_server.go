package system

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"ops-server/utils/cloud"
	cloudRequest "ops-server/utils/cloud/request"
	"time"
)

type GameServerService struct {
}

func (g *GameServerService) CreateGameServer(ctx context.Context, gameServer system.SysGameServer) (err error) {

	return global.OPS_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 获取游戏服类型信息
		var gameType system.SysGameType

		err = tx.Where("id = ?", gameServer.GameTypeId).First(&gameType).Error
		if err != nil {
			global.OPS_LOG.Error("获取游戏类型失败", zap.Error(err))
			return errors.New("获取游戏类型失败")
		}

		// 获取服务器信息
		var server system.SysAssetsServer

		err = tx.Where("id = ?", gameServer.HostId).First(&server).Error
		if err != nil {
			global.OPS_LOG.Error("获取服务器信息失败", zap.Error(err))
			return errors.New("获取服务器信息失败")
		}

		// 获取端口
		tcpPort, err := AssetsServerApp.getServerPort(server.ID, gameType.TcpPort, tx)
		if err != nil {
			global.OPS_LOG.Error("获取tpc端口失败", zap.Error(err))
			return err
		}

		httpPort, err := AssetsServerApp.getServerPort(server.ID, gameType.HttpPort, tx)
		if err != nil {
			global.OPS_LOG.Error("获取http端口失败", zap.Error(err))
			return err
		}

		grpcPort, err := AssetsServerApp.getServerPort(server.ID, gameType.GrpcPort, tx)
		if err != nil {
			global.OPS_LOG.Error("获取grpc端口失败", zap.Error(err))
			return err
		}
		// 获取vmid
		var vmid int64
		err = tx.Debug().Model(&system.SysGameServer{}).Select("IFNULL(max(vmid), 0) as max").Where("platform_id = ? and game_type_id = ?", server.PlatformId, gameType.ID).Pluck("max", &vmid).Error

		if err != nil {
			global.OPS_LOG.Error("获取vmid失败", zap.Error(err))
			return errors.New("获取vmid失败")
		}

		if vmid == 0 {
			vmid = gameType.VmidRule
		} else {
			vmid = vmid + 1
		}

		gameServer.TcpPort = tcpPort
		gameServer.HttpPort = httpPort
		gameServer.GrpcPort = grpcPort
		gameServer.Vmid = vmid

		err = tx.Create(&gameServer).Error
		if err != nil {
			global.OPS_LOG.Error("创建游戏服失败", zap.Error(err))
			return errors.New("创建游戏服失败")
		}

		// 加载关联数据
		err = tx.Where("id = ?", gameServer.ID).Preload("SysProject").Preload("Platform").Preload("GameType").Preload("Host").Preload("Redis").Preload("Mongo").Preload("Kafka").First(&gameServer).Error
		if err != nil {
			global.OPS_LOG.Error("加载游戏服关联数据失败", zap.Error(err))
			return errors.New("加载游戏服关联数据失败")
		}

		// 初始化配置文件
		gameServer.ConfigFile, err = GameTypeApp.GenerateConfigFile(gameServer)
		if err != nil {
			global.OPS_LOG.Error("初始化配置文件失败", zap.Error(err))
			return errors.New("初始化配置文件失败")
		}

		// 初始化docker-compose文件
		gameServer.ComposeFile, err = GameTypeApp.GenerateComposeFile(gameServer)
		if err != nil {
			global.OPS_LOG.Error("初始化docker-compose文件失败", zap.Error(err))
			return errors.New("初始化docker-compose文件失败")
		}

		tx.Save(&gameServer)

		return err
	})
}

func (g *GameServerService) UpdateGameServer(ctx context.Context, gameServer system.SysGameServer) (err error) {
	var old system.SysGameServer

	updateField := []string{
		"Name",
		"PlatformId",
		"GameTypeId",
		"RedisId",
		"MongoId",
		"KafkaId",
		"HostId",
	}

	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", gameServer.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	err = global.OPS_DB.WithContext(ctx).Model(&old).Select(updateField).Updates(gameServer).Error
	return
}

func (g *GameServerService) DeleteGameServer(ctx context.Context, id int) (err error) {
	//global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysGameServer{}).Error
	return global.OPS_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var gameServer system.SysGameServer
		var listenerList []system.SysAssetsListener
		if err = tx.Preload("GameType").Preload("Platform").Where("id = ?", id).First(&gameServer).Error; err != nil {
			return err
		}

		listenerName := fmt.Sprintf("%s_%s_%d", gameServer.Platform.PlatformCode, gameServer.GameType.Code, gameServer.Vmid)
		if err = tx.Preload("Lb").Preload("Lb.CloudProduce").Where("name = ?", listenerName).Find(&listenerList).Error; err != nil {
			return err
		}

		for _, listener := range listenerList {

			deleteBackendMemberParams := cloudRequest.Listener{
				AK:         listener.Lb.CloudProduce.SecretId,
				SK:         listener.Lb.CloudProduce.SecretKey,
				Region:     listener.Lb.CloudProduce.RegionId,
				ListenerId: listener.InstanceId,
			}
			if err = cloud.DeleteListener(deleteBackendMemberParams); err != nil {
				global.OPS_LOG.Error("删除监听器失败", zap.Error(err), zap.String("listenerName", listenerName), zap.String("instanceId", listener.InstanceId))
				return err
			}

			if err = tx.Unscoped().Delete(&listener).Error; err != nil {
				return err
			}
		}

		// 删除相应的游戏服目录

		sshConfig, err := task.GetSSHKey(gameServer.ProjectId, gameServer.Host.PubIp, gameServer.Host.SSHPort)
		if err != nil {
			return fmt.Errorf("获取ssh配置失败:%v", err)
		}

		sshClient, err := utils.NewSSHClient(&sshConfig)
		if err != nil {
			return fmt.Errorf("ssh连接失败:%v", err)
		}
		defer func() {
			if err := sshClient.Close(); err != nil {
				global.OPS_LOG.Error("ssh关闭失败", zap.Error(err))
			}
		}()

		gameDir := fmt.Sprintf("%s/%s/%s_%d",
			global.OPS_CONFIG.Game.GamePath,
			gameServer.Platform.PlatformCode,
			gameServer.GameType.Code,
			gameServer.Vmid,
		)

		if gameDir == "/" {
			return errors.New("游戏服目录为根目录,无法删除")
		}

		command := fmt.Sprintf("[ -d %s ] && mv  %s /tmp/", gameDir, gameDir)
		_, err = utils.ExecuteSSHCommand(sshClient, command)
		if err != nil {
			return fmt.Errorf("删除游戏服目录失败:%v", err)
		}

		return tx.Delete(&gameServer).Error
	})
}

func (g *GameServerService) GetGameServerById(ctx context.Context, id int) (result system.SysGameServer, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("Platform").First(&result, "id = ?", id).Error
	return
}

func (g *GameServerService) GetGameServerList(ctx context.Context, info request.PageInfo, server system.SysGameServer) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysGameServer{})

	var resultList []system.SysGameServer

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if server.PlatformId != 0 {
		db = db.Where("platform_id = ?", server.PlatformId)
	}

	if server.GameTypeId != 0 {
		db = db.Where("game_type_id = ?", server.GameTypeId)
	}

	if server.Status != 0 {
		db = db.Where("status = ?", server.Status)
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"

	err = db.Debug().Preload("Platform").
		Preload("Host").
		Preload("GameType", func(db *gorm.DB) *gorm.DB { return db.Select("ID,name,code") }).
		Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (g *GameServerService) GetGameServerAll(ctx context.Context) (result []system.SysGameServer, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("GameType").Find(&result).Error
	return
}

func (g *GameServerService) InstallGameServer(ctx *gin.Context, ids request.IdsReq) (jobId uuid.UUID, err error) {
	var job system.Job
	var taskList []system.JobTask
	var gameServerList []*system.SysGameServer

	// 获取需要安装的游戏服
	err = global.OPS_DB.WithContext(ctx).Preload("GameType").Preload("Host").Where("id in ?", ids.Ids).Find(&gameServerList).Error
	if err != nil {
		global.OPS_LOG.Error("获取游戏服失败", zap.Error(err))
		return
	}

	jobId = uuid.Must(uuid.NewV4())

	for index := range gameServerList {
		var t system.JobTask

		taskId := uuid.Must(uuid.NewV4())
		taskInfo, err := task.NewInstallServerTask(task.InstallServerParams{
			TaskId:       taskId,
			GameServerId: gameServerList[index].ID,
		})
		if err != nil {
			global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
			continue
		}

		t.JobId = jobId
		t.AsynqId = taskInfo.ID
		t.TaskId = taskId
		t.Status = taskInfo.State.String()
		t.HostName = gameServerList[index].Host.ServerName
		t.HostIp = gameServerList[index].Host.PubIp
		t.CreateAt = time.Now()

		if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
			global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
			continue
		}

		taskList = append(taskList, t)
		// 修改游戏服状态为安装中
		//gameServerList[index].Status = 1
		//global.OPS_DB.WithContext(ctx).Save(gameServerList[index])
	}

	job.JobId = jobId
	job.Name = "安装游戏服"
	job.Status = 1
	job.Type = "installServer"
	//job.Creator = "system"
	job.Tasks = taskList
	//job.CreateAt = time.Now()

	// 创建作业任务
	err = JobServiceApp.CreateJob(ctx, job)
	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}

	return jobId, nil
}

func (g *GameServerService) UpdateGameConfig(ctx *gin.Context, updateType int8, ids []int8) (err error) {
	var gameServerList []system.SysGameServer
	switch updateType {
	case 1:
		err = global.OPS_DB.WithContext(ctx).Where("status = 2").Preload("SysProject").Preload("Platform").Preload("GameType").Preload("Host").Preload("Redis").Preload("Mongo").Preload("Kafka").Find(&gameServerList).Error
	case 2:
		if len(ids) == 0 {
			return errors.New("选择的游戏服为空")
		}
		err = global.OPS_DB.WithContext(ctx).Where("id in ?", ids).Where("status = 2").Preload("SysProject").Preload("Platform").Preload("GameType").Preload("Host").Preload("Redis").Preload("Mongo").Preload("Kafka").Find(&gameServerList).Error
	default:
		return errors.New("更新类型错误")
	}

	if err != nil {
		return
	} else if len(gameServerList) == 0 {
		return errors.New("没有需要更新的游戏服")
	}

	for index := range gameServerList {
		// 初始化配置文件
		gameServerList[index].ConfigFile, err = GameTypeApp.GenerateConfigFile(gameServerList[index])
		if err != nil {
			global.OPS_LOG.Error("更新配置文件失败", zap.Error(err))
			return errors.New("更新配置文件失败")
		}

		// 初始化docker-compose文件
		gameServerList[index].ComposeFile, err = GameTypeApp.GenerateComposeFile(gameServerList[index])
		if err != nil {
			global.OPS_LOG.Error("更新docker-compose文件失败", zap.Error(err))
			return errors.New("更新docker-compose文件失败")
		}
		err = global.OPS_DB.WithContext(ctx).Save(&gameServerList[index]).Error
	}
	return
}

func (g *GameServerService) RsyncGameConfig(ctx *gin.Context, updateType int8, ids []int8) (jobId uuid.UUID, err error) {
	var job system.Job
	var gameServerList []system.SysGameServer
	var taskList []system.JobTask
	switch updateType {
	case 1:
		err = global.OPS_DB.WithContext(ctx).Where("status = 2").Find(&gameServerList).Error
	case 2:
		if len(ids) == 0 {
			return uuid.UUID{}, errors.New("选择的游戏服为空")
		}
		err = global.OPS_DB.WithContext(ctx).Where("id in ?", ids).Where("status = 2").Find(&gameServerList).Error
	default:
		return uuid.UUID{}, errors.New("更新类型错误")
	}

	if err != nil {
		return
	} else if len(gameServerList) == 0 {
		return uuid.UUID{}, errors.New("没有需要同步的游戏服")
	}

	// 统一获取每个主机需要更新的游戏服
	mapGameServer := make(map[uint][]uint)
	for index := range gameServerList {
		mapGameServer[gameServerList[index].HostId] = append(mapGameServer[gameServerList[index].HostId], gameServerList[index].ID)
	}

	jobId = uuid.Must(uuid.NewV4())

	for hostId, gameServerIds := range mapGameServer {
		var t system.JobTask
		var host system.SysAssetsServer

		if err = global.OPS_DB.WithContext(ctx).First(&host, "id = ?", hostId).Error; err != nil {
			global.OPS_LOG.Error("获取主机信息失败", zap.Error(err))
			continue
		}

		taskId := uuid.Must(uuid.NewV4())
		taskInfo, err := task.NewRsyncGameConfigTask(task.RsyncGameConfigParams{
			TaskId:  taskId,
			HostId:  hostId,
			GameIds: gameServerIds,
		})

		if err != nil {
			global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
			continue
		}

		t.JobId = jobId
		t.AsynqId = taskInfo.ID
		t.TaskId = taskId
		t.Status = taskInfo.State.String()
		t.HostName = host.ServerName
		t.HostIp = host.PubIp
		t.CreateAt = time.Now()

		if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
			global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
			continue
		}

		taskList = append(taskList, t)
	}

	job.JobId = jobId
	job.Name = "同步游戏服配置"
	job.Status = 1
	job.Type = task.RsyncGameConfigTypeName
	//job.Creator = "system"
	job.Tasks = taskList

	// 创建作业任务
	err = JobServiceApp.CreateJob(ctx, job)
	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}

	return
}

func (g *GameServerService) ExecGameTask(ctx *gin.Context, taskType int8, ids []uint) (jobId uuid.UUID, err error) {
	var gameServerList []system.SysGameServer
	var taskTypeName, jobName string
	var taskList []system.JobTask
	var job system.Job

	if len(ids) == 0 {
		return uuid.UUID{}, errors.New("选择的游戏服为空")
	}

	err = global.OPS_DB.WithContext(ctx).Where("id in ?", ids).Where("status = 2").Preload("Platform").Preload("GameType").Preload("Host").Preload("Redis").Preload("Mongo").Preload("Kafka").Find(&gameServerList).Error

	if err != nil {
		return uuid.UUID{}, errors.New("获取游戏服信息失败")
	}

	switch taskType {
	case 1:
		// 开启游戏服
		taskTypeName = task.StartGameTypeName
		jobName = "开启游戏服"
	case 2:
		// 关闭游戏服
		taskTypeName = task.StopGameTypeName
		jobName = "关闭游戏服"
	case 3:
		// 更新游戏服配置文件
		taskTypeName = task.UpdateGameConfigTypeName
		jobName = "更新游戏服配置文件"
	case 4:
		// 同步游戏服配置文件
		taskTypeName = task.RsyncGameConfigTypeName
		jobName = "同步游戏服配置文件"
	case 5:
		// 安装游戏服
		taskTypeName = task.InstallGameServerTypeName
		jobName = "安装游戏服"

	default:
		return uuid.UUID{}, errors.New("更新类型错误")
	}

	mapGameServer := make(map[uint][]uint)

	// 整理游戏服信息
	for index := range gameServerList {
		if _, ok := mapGameServer[gameServerList[index].HostId]; !ok {
			mapGameServer[gameServerList[index].HostId] = make([]uint, 0)
		}
		mapGameServer[gameServerList[index].HostId] = append(mapGameServer[gameServerList[index].HostId], gameServerList[index].ID)
	}

	jobId = uuid.Must(uuid.NewV4())

	for hostId, gameServerIds := range mapGameServer {
		var t system.JobTask
		var host system.SysAssetsServer

		if err = global.OPS_DB.WithContext(ctx).First(&host, "id = ?", hostId).Error; err != nil {
			global.OPS_LOG.Error("获取主机信息失败", zap.Error(err))
			continue
		}

		taskId := uuid.Must(uuid.NewV4())
		taskInfo, err := task.NewGameTask(taskTypeName, task.GameTaskParams{
			TaskId:        taskId,
			HostId:        hostId,
			GameServerIds: gameServerIds,
			ProjectId:     host.ProjectId,
		})

		if err != nil {
			global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
			continue
		}

		t.JobId = jobId
		t.AsynqId = taskInfo.ID
		t.TaskId = taskId
		t.Status = taskInfo.State.String()
		t.HostName = host.ServerName
		t.HostIp = host.PubIp
		t.CreateAt = time.Now()

		if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
			global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
			continue
		}

		taskList = append(taskList, t)
	}
	job.JobId = jobId
	job.Name = jobName
	job.Status = 1
	job.Type = taskTypeName
	job.Tasks = taskList

	// 创建作业任务
	err = JobServiceApp.CreateJob(ctx, job)
	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}
	return
}
