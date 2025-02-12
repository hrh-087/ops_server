package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
	"path/filepath"
	"time"
)

type RsyncGameConfigParams struct {
	TaskId  uuid.UUID
	HostId  uint
	GameIds []uint
}

func NewRsyncGameConfigTask(params RsyncGameConfigParams) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(RsyncGameConfigTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

func HandleRsyncGameConfigTask(ctx context.Context, t *asynq.Task) (err error) {
	var params RsyncGameConfigParams
	var resultList []string
	var host system.SysAssetsServer
	var gameServerList []system.SysGameServer

	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	if err = global.OPS_DB.Preload("Platform").Preload("GameType").Preload("Host").Where("id = ?", params.HostId).First(&host).Error; err != nil {
		return fmt.Errorf("获取负载均衡信息失败:%v", err)
	}

	if err = global.OPS_DB.Preload("Platform").Preload("GameType").Preload("Host").Where("id in (?)", params.GameIds).Find(&gameServerList).Error; err != nil {
		return fmt.Errorf("获取游戏服信息失败:%v", err)
	}

	defer func() {
		if err != nil {
			global.OPS_LOG.Error("同步游戏服配置失败", zap.Error(err))
			// 写入任务结果
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	sshConfig, err := GetSSHKey(host.ProjectId, host.PubIp, host.SSHPort)
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

	for _, gameServer := range gameServerList {
		resultList = append(resultList, fmt.Sprintf("开始同步%s_%d配置文件...", gameServer.GameType.Code, gameServer.Vmid))
		// 定义游戏服目录
		gameServerDir := fmt.Sprintf("%s/%s/%s_%d",
			global.OPS_CONFIG.Game.GamePath,
			gameServer.Platform.PlatformCode,
			gameServer.GameType.Code,
			gameServer.Vmid,
		)
		// 本地保存目录
		localGameServerDir := fmt.Sprintf("%s/%s/%s_%d",
			filepath.Join(global.OPS_CONFIG.Local.Path, "gameConfig", time.Now().Format("2006-01-02")),
			gameServer.Platform.PlatformCode,
			gameServer.GameType.Code,
			gameServer.Vmid,
		)

		configFilePath, err := utils.CreateFile(localGameServerDir, "application.yaml", gameServer.ConfigFile)
		if err != nil {
			return fmt.Errorf("创建配置文件失败:%v", err)
		}

		composeFilePath, err := utils.CreateFile(localGameServerDir, "docker-compose.yaml", gameServer.ComposeFile)
		if err != nil {
			return fmt.Errorf("创建docker-compose文件失败:%v", err)
		}

		err = utils.UploadFile(sshClient, configFilePath, fmt.Sprintf("%s/data/application.yaml", gameServerDir))
		if err != nil {
			return fmt.Errorf("上传配置文件失败:%v", err)
		}

		err = utils.UploadFile(sshClient, composeFilePath, fmt.Sprintf("%s/docker-compose.yaml", gameServerDir))
		if err != nil {
			return fmt.Errorf("上传docker-compose文件失败:%v", err)
		}
	}
	resultList = append(resultList, "同步成功...")
	WriteTaskResult(t, resultList)

	return
}
