package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
)

type InstallServerParams struct {
	TaskId       uuid.UUID
	GameServerId uint
}

func NewInstallServerTask(params InstallServerParams) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(InstallServerTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

func HandleInstallServerTask(ctx context.Context, t *asynq.Task) error {
	var gameServer system.SysGameServer
	//var task system.JobTask
	var resultList []string
	var params InstallServerParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	if err := global.OPS_DB.Preload("Platform").Preload("GameType").Preload("Host").Where("id = ?", params.GameServerId).First(&gameServer).Error; err != nil {
		resultList = append(resultList, "获取游戏服信息失败")
		global.OPS_LOG.Error("获取游戏服信息失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	if gameServer.Status == 1 || gameServer.Status == 2 {
		resultList = append(resultList, "该服务已安装或正在安装")
		WriteTaskResult(t, resultList)
		return errors.New("该服务已安装或正在安装")
	}
	// 修改游戏服状态
	gameServer.Status = 1
	global.OPS_DB.Save(&gameServer)

	sshConfig, err := GetSSHKey(gameServer.ProjectId, gameServer.Host.PubIp, gameServer.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "获取ssh配置失败")
		global.OPS_LOG.Error("获取ssh配置失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := utils.NewSSHClient(&sshConfig)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// 定义游戏服目录
	gameServerDir := fmt.Sprintf("%s/%s/%s_%d",
		global.OPS_CONFIG.Game.GamePath,
		gameServer.Platform.PlatformCode,
		gameServer.GameType.Code,
		gameServer.Vmid,
	)

	command := fmt.Sprintf("mkdir -p %s/data/hotswap", gameServerDir)
	_, err = utils.ExecuteSSHCommand(sshClient, command)
	if err != nil {
		resultList = append(resultList, "创建游戏服目录失败")
		global.OPS_LOG.Error("创建游戏服目录失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	configFilePath, err := utils.CreateFile(gameServerDir, "application.yaml", gameServer.ConfigFile)
	if err != nil {
		resultList = append(resultList, "创建配置文件失败")
		global.OPS_LOG.Error("创建配置文件失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	composeFilePath, err := utils.CreateFile(gameServerDir, "docker-compose.yaml", gameServer.ComposeFile)
	if err != nil {
		resultList = append(resultList, "创建docker-compose文件失败")
		global.OPS_LOG.Error("创建docker-compose文件失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	err = utils.UploadFile(sshClient, configFilePath, fmt.Sprintf("%s/data/application.yaml", gameServerDir))
	if err != nil {
		resultList = append(resultList, "上传配置文件失败")
		global.OPS_LOG.Error("上传配置文件失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	err = utils.UploadFile(sshClient, composeFilePath, fmt.Sprintf("%s/docker-compose.yaml", gameServerDir))
	if err != nil {
		resultList = append(resultList, "上传docker-compose文件失败")
		global.OPS_LOG.Error("上传docker-compose文件失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, "安装成功")
	WriteTaskResult(t, resultList)

	gameServer.Status = 2
	global.OPS_DB.Save(&gameServer)

	return err
}
