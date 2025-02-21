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

// 关闭游戏服
func HandleStartGame(ctx context.Context, t *asynq.Task) (err error) {
	var resultList []string
	var params GameTaskParams
	var output string

	defer func() {
		if err != nil {
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	if err = json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	output, err = StartGameServer(params.ProjectId, params.HostId, params.GameServerIds)
	resultList = append(resultList, output)
	if err != nil {
		global.OPS_LOG.Error("开启游戏服失败", zap.Error(err))
		return err
	}
	WriteTaskResult(t, resultList)
	return nil

}

func StartGameServer(projectId uint, hostId uint, gameServerIds []uint) (output string, err error) {
	var host system.SysAssetsServer
	var gameServerList []system.SysGameServer
	var outputList []string

	err = global.OPS_DB.Where("id = ?", hostId).First(&host).Error
	if err != nil {
		return "", err
	}

	// 当gameServerIds为空时，获取该主机上所有游戏服
	if len(gameServerIds) == 0 {
		if err = global.OPS_DB.Where("host_id = ?", hostId).Preload("Platform").Preload("GameType").Find(&gameServerList).Error; err != nil {
			return "", err
		}
	} else {
		if err = global.OPS_DB.Where("id in ?", gameServerIds).Preload("Platform").Preload("GameType").Find(&gameServerList).Error; err != nil {
			return "", err
		}
	}
	sshClient, err := GetSSHConn(projectId, host.PubIp, host.SSHPort)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := sshClient.Close(); err != nil {
			global.OPS_LOG.Error("ssh连接关闭失败", zap.Error(err))
		}
	}()

	for _, gameServer := range gameServerList {
		// 脚本命令 渠道名称 游戏服目录
		command := fmt.Sprintf("bash %s %s %s_%d", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "start_game.sh"), gameServer.Platform.PlatformCode, gameServer.GameType.Code, gameServer.Vmid)
		outputList = append(outputList, command)
		// windows开发端路径替换
		command = strings.ReplaceAll(command, "\\", "/")
		// 执行关闭游戏服命令
		output, err = utils.ExecuteSSHCommand(sshClient, command)
		outputList = append(outputList, output)
		if err != nil {
			return strings.Join(outputList, "\n"), err
		}
	}

	return strings.Join(outputList, "\n"), err
}
