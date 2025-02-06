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
		err = tx.Where("id = ?", gameServer.ID).Preload("Platform").Preload("GameType").Preload("Host").Preload("Redis").Preload("Mongo").Preload("Kafka").First(&gameServer).Error
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
		if err := tx.Preload("GameType").Preload("Platform").Where("id = ?", id).First(&gameServer).Error; err != nil {
			return err
		}

		listenerName := fmt.Sprintf("%s_%s_%d", gameServer.Platform.PlatformCode, gameServer.GameType.Code, gameServer.Vmid)
		if err := tx.Preload("Lb").Preload("Lb.CloudProduce").Where("name = ?", listenerName).Find(&listenerList).Error; err != nil {
			return err
		}

		for _, listener := range listenerList {

			deleteBackendMemberParams := cloudRequest.Listener{
				AK:         listener.Lb.CloudProduce.SecretId,
				SK:         listener.Lb.CloudProduce.SecretKey,
				Region:     listener.Lb.CloudProduce.RegionId,
				ListenerId: listener.InstanceId,
			}
			if err := cloud.DeleteListener(deleteBackendMemberParams); err != nil {
				global.OPS_LOG.Error("删除监听器失败", zap.Error(err), zap.String("listenerName", listenerName), zap.String("instanceId", listener.InstanceId))
				return err
			}

			if err := tx.Unscoped().Delete(&listener).Error; err != nil {
				return err
			}
		}

		return tx.Delete(&gameServer).Error
	})
}

func (g *GameServerService) GetGameServerById(ctx context.Context, id int) (result system.SysGameServer, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("Platform").First(&result, "id = ?", id).Error
	return
}

func (g *GameServerService) GetGameServerList(ctx context.Context, info request.PageInfo, server request.NameAndPlatformSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysGameServer{})

	var resultList []system.SysGameServer

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if server.Name != "" {
		db = db.Where("name like ?", "%"+server.Name+"%")
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"

	err = db.Preload("Platform").
		Preload("Host").
		Preload("GameType", func(db *gorm.DB) *gorm.DB { return db.Select("ID,name") }).
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
