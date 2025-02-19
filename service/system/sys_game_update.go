package system

import (
	"encoding/json"
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
	sysReq "ops-server/model/system/request"
	"path/filepath"
	"strings"
	"time"
)

type GameUpdateService struct {
}

// 创建正常更新步骤
func (s *GameUpdateService) createNormalUpdateStep(step int8, version string) (updateParams sysReq.GameUpdateTaskParams) {
	switch step {
	case 1:
		updateParams.StepName = "拉取游戏服镜像"
		updateParams.TaskTypeName = task.UpdateGameImageTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())

	case 2:
		updateParams.StepName = "关闭游戏服"
		updateParams.TaskTypeName = task.StopGameTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())

	case 3:
		updateParams.StepName = "更新游戏服配置"
		updateParams.TaskTypeName = task.RsyncGameJsonConfigTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())

	case 4:
		updateParams.StepName = "开启游戏服"
		updateParams.TaskTypeName = task.StartGameTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())
	case 5:
		updateParams.StepName = "检查游戏版本号"
		updateParams.TaskTypeName = task.CheckGameVersionTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())
	}
	return
}

func (s *GameUpdateService) createHotUpdateStep(step int8, hotParams sysReq.HotUpdateParams) (updateParams sysReq.GameUpdateTaskParams) {

	//params, _ := json.Marshal(hotParams)

	switch step {
	case 1:
		updateParams.StepName = "解压热更包"
		updateParams.TaskTypeName = task.HotGameUnzipFileTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())
		updateParams.Command = "unzip -o %s -d /tmp/%s"
	case 2:
		updateParams.StepName = "同步相应服务器"
		updateParams.TaskTypeName = task.HotGameRsyncHostTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())
		updateParams.Command = fmt.Sprintf("sh %s", filepath.Join(global.OPS_CONFIG.Game.GameScriptPath, "hot_game_rsync_host.sh"))
	case 3:
		updateParams.StepName = "同步到相应游戏服"
		updateParams.TaskTypeName = task.HotGameRsyncServerTypeName
		updateParams.JobId = uuid.Must(uuid.NewV4())

		if hotParams.ServerType == 1 {
			updateParams.Command = fmt.Sprintf("sh %s", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "hot_game_rsync_server.sh game "))
		} else {
			updateParams.Command = fmt.Sprintf("sh %s", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "hot_game_rsync_server.sh game_type "))
		}
	}
	return
}

func (s *GameUpdateService) CreateGameUpdate(ctx *gin.Context, gameUpdate system.GameUpdate, hotParams sysReq.HotUpdateParams) (err error) {
	return global.OPS_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch gameUpdate.UpdateType {
		case 1:
			updateParams := make(map[int]sysReq.GameUpdateTaskParams)
			// 定义正常更新总步骤数
			totalStep := 5
			for i := 1; i <= totalStep; i++ {
				updateParams[i] = s.createNormalUpdateStep(int8(i), gameUpdate.GameVersion)
			}

			params, err := json.Marshal(updateParams)
			if err != nil {
				return err
			}

			gameUpdate.TotalStep = int8(totalStep)
			gameUpdate.UpdateParams = string(params)
			gameUpdate.Step = 1
			gameUpdate.StepName = updateParams[1].StepName

		case 2:
			updateParams := make(map[int]sysReq.GameUpdateTaskParams)
			// 定义正常更新总步骤数
			totalStep := 3
			// 获取热更步骤参数
			for i := 1; i <= totalStep; i++ {
				updateParams[i] = s.createHotUpdateStep(int8(i), hotParams)
			}

			hotServerList, err := json.Marshal(hotParams.ServerList)
			if err != nil {
				return err
			}
			params, err := json.Marshal(updateParams)
			if err != nil {
				return err
			}

			gameUpdate.ServerType = hotParams.ServerType
			gameUpdate.ServerList = string(hotServerList)
			gameUpdate.TotalStep = int8(totalStep)
			gameUpdate.UpdateParams = string(params)
			gameUpdate.Step = 1
			gameUpdate.StepName = updateParams[1].StepName

		case 3:
			fmt.Println("热更游戏配置")
		default:
			return errors.New("未知类型")
		}
		err = tx.Create(&gameUpdate).Error
		return err
	})
}

func (s *GameUpdateService) UpdateGameUpdate(ctx *gin.Context, gameUpdate system.GameUpdate, hotParams sysReq.HotUpdateParams) (err error) {
	return
}

func (s *GameUpdateService) DeleteGameUpdate(ctx *gin.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Unscoped().Delete(&system.GameUpdate{}, "id = ?", id).Error
}

func (s *GameUpdateService) GetGameUpdateList(ctx *gin.Context, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.GameUpdate{})

	var resultList []system.GameUpdate

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error

	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"

	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (s *GameUpdateService) GetGameUpdateById(ctx *gin.Context, id int) (result system.GameUpdate, err error) {
	err = global.OPS_DB.WithContext(ctx).First(&result, "id = ?", id).Error
	if err != nil {
		return result, err
	}

	updateParams := make(map[int8]sysReq.GameUpdateTaskParams)
	if err = json.Unmarshal([]byte(result.UpdateParams), &updateParams); err != nil {
		return result, err
	}

	var job system.Job
	err = global.OPS_DB.WithContext(ctx).First(&job, "job_id = ?", updateParams[result.Step].JobId).Error
	if err == gorm.ErrRecordNotFound {
		return result, nil
	} else if err != nil {
		return result, err
	}

	if job.Status == 2 {
		newStep := result.Step + 1
		if _, exits := updateParams[newStep]; !exits {
			result.Step = 0
			result.StepName = "已完成"
		} else {
			result.Step = newStep
			result.StepName = updateParams[newStep].StepName
		}

		err = global.OPS_DB.WithContext(ctx).Save(&result).Error
		return result, err
	}
	return
}

func (s *GameUpdateService) ExecUpdateTask(ctx *gin.Context, id int) (jobId uuid.UUID, err error) {
	var job system.Job
	var taskList []system.JobTask
	//var hostList []system.SysAssetsServer
	var gameUpdate system.GameUpdate

	if err = global.OPS_DB.WithContext(ctx).Preload("SysProject").First(&gameUpdate, "id = ?", id).Error; err != nil {
		return
	}

	updateParams := make(map[int8]sysReq.GameUpdateTaskParams)
	if err = json.Unmarshal([]byte(gameUpdate.UpdateParams), &updateParams); err != nil {
		return
	}

	if _, exists := updateParams[gameUpdate.Step]; !exists {
		return jobId, errors.New("步骤出错,请联系相关人员排查")
	}

	// 验证jobid是否已执行过，执行后续重新替换jobid
	if errors.Is(global.OPS_DB.WithContext(ctx).First(&job, "job_id = ?", updateParams[gameUpdate.Step].JobId).Error, gorm.ErrRecordNotFound) {
		jobId = updateParams[gameUpdate.Step].JobId
	} else {
		jobId = uuid.Must(uuid.NewV4())
		updateParams[gameUpdate.Step] = sysReq.GameUpdateTaskParams{
			JobId:        jobId,
			TaskTypeName: updateParams[gameUpdate.Step].TaskTypeName,
			Command:      updateParams[gameUpdate.Step].Command,
			StepName:     updateParams[gameUpdate.Step].StepName,
		}
		params, err := json.Marshal(updateParams)
		if err != nil {
			return jobId, err
		}
		gameUpdate.UpdateParams = string(params)
		global.OPS_DB.WithContext(ctx).Save(&gameUpdate)
	}

	switch gameUpdate.UpdateType {
	// 正常更新
	case 1:
		var hostIdList []int
		// 根据不同的步骤获取主机列表
		switch gameUpdate.Step {
		case 1, 2, 4:
			err = global.OPS_DB.WithContext(ctx).Model(&system.SysGameServer{}).Where("status = ?", 2).Group("host_id").Pluck("host_id", &hostIdList).Error
			if err != nil {
				return
			}
		case 3, 5:
			err = global.OPS_DB.Debug().WithContext(ctx).Model(&system.SysAssetsServer{}).Where("status = ? and server_type = ?", 1, 3).Limit(1).Pluck("id", &hostIdList).Error
			if err == gorm.ErrRecordNotFound {
				return jobId, errors.New("未添加后台服务器")
			}

		}

		if len(hostIdList) == 0 {
			return jobId, errors.New("未获取到匹配主机")
		}

		//if gameUpdate.Step == 3 {
		//	command = fmt.Sprintf("%s %s %s %s",
		//		updateParams[gameUpdate.Step].Command,
		//		gameUpdate.SysProject.ConfigDir,
		//		//global.OPS_CONFIG.Game.GameConfigDir, // 临时写死
		//		global.OPS_CONFIG.Game.RemoteConfigDir,
		//		strings.Join(allHostIpList, ","),
		//	)
		//} else if gameUpdate.Step == 5 {
		//	command = fmt.Sprintf("%s %s %s", updateParams[gameUpdate.Step].Command, gameUpdate.GameVersion, strings.Join(allHostIpList, ","))
		//} else {
		//	command = updateParams[gameUpdate.Step].Command
		//}

		// 根据步骤参数添加到任务列表中
		for _, hostId := range hostIdList {
			var t system.JobTask
			var host system.SysAssetsServer
			var taskId uuid.UUID
			if err = global.OPS_DB.WithContext(ctx).First(&host, "id = ?", hostId).Error; err != nil {
				global.OPS_LOG.Error("获取主机信息失败", zap.Error(err))
				continue
			}
			taskId = uuid.Must(uuid.NewV4())
			taskInfo, err := task.NewGameTask(updateParams[gameUpdate.Step].TaskTypeName, task.GameTaskParams{
				HostId:    host.ID,
				TaskId:    taskId,
				ProjectId: gameUpdate.ProjectId,
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

	//热更
	case 2:
		var hostIdList []int
		var gameServerList []*system.SysGameServer

		// 获取执行主机
		switch gameUpdate.Step {
		case 1, 2:
			err = global.OPS_DB.WithContext(ctx).Model(&system.SysAssetsServer{}).Where("status = ? and server_type = ?", 1, 3).Limit(1).Pluck("id", &hostIdList).Error
			if err == gorm.ErrRecordNotFound {
				return jobId, errors.New("未添加后台服务器")
			}
		case 3:
			var gameTypeList []int
			if err = json.Unmarshal([]byte(gameUpdate.ServerList), &gameTypeList); err != nil {
				return
			}

			switch gameUpdate.ServerType {
			case 1:
				// 游戏服
				err = global.OPS_DB.Debug().WithContext(ctx).Model(&system.SysGameServer{}).Preload("GameType").Preload("Host").Where("status = ? and id in ?", 2, gameTypeList).Find(&gameServerList).Error
				if err != nil {
					return
				}
			case 2:
				// 游戏服类型
				err = global.OPS_DB.Debug().WithContext(ctx).Model(&system.SysGameServer{}).Preload("GameType").Preload("Host").Where("status = ? and game_type_id in ?", 2, gameTypeList).Find(&gameServerList).Group("game_type_id,host_id").Error
				if err != nil {
					return
				}
			}
		}

		// 获取热更文件信息，替换文件路径
		hotFilePath := strings.ReplaceAll(gameUpdate.HotFile, "resource", global.OPS_CONFIG.Game.HotFileDir)
		hotFileName := strings.Split(filepath.Base(hotFilePath), ".")[0]
		if hotFileName == "" {
			return jobId, errors.New("热更文件不能为空")
		}

		if len(hostIdList) > 0 {
			var serverList []system.SysGameServer
			var commandParams []string
			var command string

			if gameUpdate.Step == 2 {
				var gameTypeList []int
				if err = json.Unmarshal([]byte(gameUpdate.ServerList), &gameTypeList); err != nil {
					return
				}

				switch gameUpdate.ServerType {
				case 1:
					// 游戏服
					err = global.OPS_DB.Debug().WithContext(ctx).Model(&system.SysGameServer{}).Preload("Host").Where("status = ? and id in ?", 2, gameTypeList).Find(&serverList).Error
					if err != nil {
						return
					}
				case 2:
					// 游戏服类型
					err = global.OPS_DB.Debug().WithContext(ctx).Model(&system.SysGameServer{}).Preload("Host").Where("status = ? and game_type_id in ?", 2, gameTypeList).Find(&serverList).Group("game_type_id,host_id").Error
					if err != nil {
						return
					}
				}

				for _, server := range serverList {
					commandParams = append(commandParams, server.Host.PubIp)
				}
				// 拼接命令
				command = fmt.Sprintf("%s /tmp/%s %s", updateParams[gameUpdate.Step].Command, hotFileName, strings.Join(commandParams, ","))

			} else if gameUpdate.Step == 1 {
				command = fmt.Sprintf(updateParams[gameUpdate.Step].Command, hotFilePath, hotFileName)
			} else {
				command = updateParams[gameUpdate.Step].Command
			}

			for _, hostId := range hostIdList {
				var t system.JobTask
				var host system.SysAssetsServer
				taskId := uuid.Must(uuid.NewV4())
				if err = global.OPS_DB.WithContext(ctx).First(&host, "id = ?", hostId).Error; err != nil {
					global.OPS_LOG.Error("获取主机信息失败", zap.Error(err))
					continue
				}

				taskInfo, err := task.NewUpdateGameTask(updateParams[gameUpdate.Step].TaskTypeName, task.NormalUpdateGameParams{
					Host:    host,
					TaskId:  taskId,
					Command: command,
					Params:  updateParams[gameUpdate.Step].Params,
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
		} else if len(gameServerList) > 0 {
			for _, gameServer := range gameServerList {
				var t system.JobTask
				var command string

				taskId := uuid.Must(uuid.NewV4())
				// 根据热更类型拼接命令
				if gameUpdate.ServerType == 1 {
					command = fmt.Sprintf("%s %s_%d /tmp/%s/", updateParams[gameUpdate.Step].Command, gameServer.GameType.Code, gameServer.Vmid, hotFileName)
				} else if gameUpdate.ServerType == 2 {
					command = fmt.Sprintf("%s %s /tmp/%s/", updateParams[gameUpdate.Step].Command, gameServer.GameType.Code, hotFileName)
				} else {
					return jobId, errors.New("无法识别游戏服热更类型")
				}

				taskInfo, err := task.NewUpdateGameTask(updateParams[gameUpdate.Step].TaskTypeName, task.NormalUpdateGameParams{
					Host:    gameServer.Host,
					TaskId:  taskId,
					Command: command,
					Params:  updateParams[gameUpdate.Step].Params,
				})
				if err != nil {
					global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
					continue
				}

				t.JobId = jobId
				t.AsynqId = taskInfo.ID
				t.TaskId = taskId
				t.Status = taskInfo.State.String()
				t.HostName = gameServer.Host.ServerName
				t.HostIp = gameServer.Host.PubIp
				t.CreateAt = time.Now()

				if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
					global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
					continue
				}

				taskList = append(taskList, t)
			}
		}
	case 3:
		return uuid.UUID{}, errors.New("暂不支持更新配置文件")
	default:
		return uuid.UUID{}, errors.New("无法识别更新类型,请联系相关成员排查")
	}

	// 创建作业任务
	job.JobId = jobId
	job.Name = updateParams[gameUpdate.Step].StepName
	job.Status = 1
	job.Type = updateParams[gameUpdate.Step].TaskTypeName
	job.Tasks = taskList

	if err = JobServiceApp.CreateJob(ctx, job); err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}
	return jobId, err

}
