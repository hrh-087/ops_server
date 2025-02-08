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
	"strings"
)

//type UpdateGameImageParams struct {
//	TaskId uuid.UUID
//}
//
//type StopGameParams struct {
//	TaskId uuid.UUID
//}

type NormalUpdateGameParams struct {
	TaskId  uuid.UUID
	Host    system.SysAssetsServer
	Command string
	Params  string
}

func NewUpdateGameTask(taskTypeName string, params interface{}) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(taskTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
}

// 更新游戏服镜像
func HandleUpdateGameImage(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)

	return nil
}

// 关闭游戏服
func HandleStopGame(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 更新游戏服配置
func HandleUpdateGameJsonData(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 开启游戏服
func HandleStartGame(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 检查版本号
func HandleCheckGameVersion(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 热更
// 解压安装包
func HandleHotGameUnzipFile(ctx context.Context, t *asynq.Task) error {

	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 同步热更文件到相应服务器
func HandleHotGameRsyncHost(ctx context.Context, t *asynq.Task) error {

	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}

// 同步热更文件到相应游戏服
func HandleHotGameRsyncServer(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	sshClient, err := GetSSHConn(params.Host.ProjectId, params.Host.PubIp, params.Host.SSHPort)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	// windows开发端路径替换
	params.Command = strings.ReplaceAll(params.Command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", params.Command))
	output, err := utils.ExecuteSSHCommand(sshClient, params.Command)
	if err != nil {
		resultList = append(resultList, output)
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}
