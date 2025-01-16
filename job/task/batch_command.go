package task

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"ops-server/utils"
)

type BatchCommand struct {
	TaskId  uuid.UUID
	Command string
	Host    system.SysAssetsServer
}

func NewBatchCommand(params BatchCommand) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(BatchCommandTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

func HandleBatchCommand(ctx context.Context, t *asynq.Task) error {
	var params BatchCommand
	var resultList []string

	if err := json.Unmarshal(t.Payload(), &params); err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshConfig, err := GetSSHKey(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
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

	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	resultList = append(resultList, output)
	if err != nil {
		global.OPS_LOG.Error("命令执行失败:", zap.Error(err))
		return err
	}

	WriteTaskResult(t, resultList)
	return nil
}
