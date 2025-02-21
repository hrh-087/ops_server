package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
	"path/filepath"
	"strings"
)

// 更新游戏服镜像
func HandleUpdateGameImage(ctx context.Context, t *asynq.Task) (err error) {
	var resultList []string
	var params GameTaskParams

	defer func() {
		if err != nil {
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	output, err := UpdateGameImage(params.ProjectId, params.HostId)

	// 不管执行成功还是失败都要写入结果
	resultList = append(resultList, output)
	if err != nil {
		global.OPS_LOG.Error("更新游戏镜像失败", zap.Error(err))
		return err
	}
	WriteTaskResult(t, resultList)

	return nil
}

func UpdateGameImage(projectId uint, hostId uint) (output string, err error) {

	var host system.SysAssetsServer

	// 拼接命令
	command := fmt.Sprintf("bash %s", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "pull_game_image.sh"))

	err = global.OPS_DB.Where("id = ?", hostId).First(&host).Error
	if err != nil {
		return "", err
	}

	sshClient, err := GetSSHConn(projectId, host.PubIp, host.SSHPort)
	if err != nil {
		return "", err
	}
	defer func() {
		err := sshClient.Close()
		if err != nil {
			global.OPS_LOG.Error("ssh连接关闭失败", zap.Error(err))
		}
	}()

	// windows开发端路径替换
	command = strings.ReplaceAll(command, "\\", "/")

	global.OPS_LOG.Info("拉取游戏服镜像命令", zap.String("command", command))

	return utils.ExecuteSSHCommand(sshClient, command)
}
