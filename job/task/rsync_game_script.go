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
	"strings"
)

type RsyncGameScriptParams struct {
	TaskId       uuid.UUID
	TaskManageId uint
	HostId       uint
}

func HandleRsyncGameScriptTask(ctx context.Context, t *asynq.Task) (err error) {
	var params CommonTaskParams
	var resultList []string
	var host system.SysAssetsServer
	var ipList []string

	defer func() {
		if err != nil {
			global.OPS_LOG.Error("同步游戏服脚本失败", zap.Error(err))
			// 写入任务结果
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	if err = json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	if err = global.OPS_DB.Where("id = ?", params.HostId).First(&host).Error; err != nil {
		return err
	}

	if err = global.OPS_DB.Model(&system.SysAssetsServer{}).Where("server_type != 3 and status = 1").Pluck("pub_ip", &ipList).Error; err != nil {
		return err
	}

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

	command := fmt.Sprintf("sh %s %s %s",
		filepath.Join(global.OPS_CONFIG.Game.GameScriptPath, "rsync_script.sh"),
		global.OPS_CONFIG.Game.GameScriptAutoPath,
		strings.Join(ipList, ","),
	)

	output, err := utils.ExecuteSSHCommand(sshClient, command)
	resultList = append(resultList, output)
	if err != nil {
		return err
	}

	WriteTaskResult(t, resultList)

	return err

}
