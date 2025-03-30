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

func HandleRsyncGameJsonConfig(ctx context.Context, t *asynq.Task) (err error) {
	var params GameTaskParams
	var resultList []string

	defer func() {
		if err != nil {
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	output, err := RsyncGameJsonConfig(params.ProjectId, 0)

	// 不管执行成功还是失败都要写入结果
	resultList = append(resultList, output)
	if err != nil {
		global.OPS_LOG.Error("同步配置文件失败", zap.Error(err))
		return err
	}
	resultList = append(resultList, "同步配置文件成功")
	WriteTaskResult(t, resultList)

	return err
}

func RsyncGameJsonConfig(projectId uint, hostId uint) (output string, err error) {
	var project system.SysProject
	var hostIpList []string

	// 获取项目信息
	err = global.OPS_DB.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		return "", err
	}

	if hostId == 0 {
		// 获取游戏服主机
		//err = global.OPS_DB.Model(&system.SysAssetsServer{}).Where("project_id = ? and server_type = 1", projectId).Pluck("pub_ip", &hostIpList).Error
		err = global.OPS_DB.Model(&system.SysGameServer{}).Select("host.pub_ip").Joins("join sys_assets_servers as host on host.id = sys_game_servers.host_id").Where("host.project_id = ? and host.server_type = 1", projectId).Group("host_id").Find(&hostIpList).Error
		if err != nil {
			return "", err
		}
	} else {
		err = global.OPS_DB.Model(&system.SysAssetsServer{}).Where("id = ? and project_id = ?", hostId, projectId).Pluck("pub_ip", &hostIpList).Error
		if err != nil {
			return "", err
		}
	}

	sshClient, err := GetSSHConn(projectId, global.OPS_CONFIG.Ops.Host, global.OPS_CONFIG.Ops.Port)
	if err != nil {
		return "", err
	}
	defer func() {
		err := sshClient.Close()
		if err != nil {
			global.OPS_LOG.Error("ssh连接关闭失败", zap.Error(err))
		}
	}()

	command := fmt.Sprintf(
		"bash %s %s %s %s",
		filepath.Join(global.OPS_CONFIG.Game.GameScriptPath, "game_sync_config.sh"),
		project.ConfigDir,
		global.OPS_CONFIG.Game.RemoteConfigDir,
		strings.Join(hostIpList, ","),
	)
	// windows开发端路径替换
	command = strings.ReplaceAll(command, "\\", "/")

	global.OPS_LOG.Info("同步配置文件命令", zap.String("command", command))

	return utils.ExecuteSSHCommand(sshClient, command)
}
