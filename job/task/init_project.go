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

type InitProjectParams struct {
	TaskId  uuid.UUID
	Project system.SysProject
}

func NewInitProjectTask(params InitProjectParams) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(InitProjectTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

func HandleInitProject(ctx context.Context, t *asynq.Task) (err error) {
	var params InitProjectParams
	var resultList []string

	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			global.OPS_LOG.Error("初始化项目失败", zap.Error(err))
			// 写入任务结果
			resultList = append(resultList, err.Error())
			WriteTaskResult(t, resultList)
		}
	}()

	sshConfig, err := GetSSHKey(params.Project.ID, global.OPS_CONFIG.Ops.Host, global.OPS_CONFIG.Ops.Port)
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

	command := fmt.Sprintf("bash %s %s %s",
		filepath.Join(global.OPS_CONFIG.Game.GameScriptPath, "init_project.sh"),
		params.Project.ConfigDir,
		params.Project.SvnUrl,
	)

	command = strings.ReplaceAll(command, "\\", "/")

	output, err := utils.ExecuteSSHCommand(sshClient, command)
	resultList = append(resultList, output)
	if err != nil {
		return err
	}

	params.Project.Status = 2
	if err := global.OPS_DB.Save(&params.Project).Error; err != nil {
		return err
	}

	resultList = append(resultList, "初始化成功...")
	WriteTaskResult(t, resultList)

	return err
}
