package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/utils"
	"path/filepath"
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
	TaskId uuid.UUID
	//Host        system.SysAssetsServer
	ServerType  int8
	ProjectId   uint
	HotFilePath string
	HotFileName string
	IpList      string
	GameType    string
	GameVmid    int64
}

func NewUpdateGameTask(taskTypeName string, params interface{}) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(taskTypeName, payload)
	return global.AsynqClinet.Enqueue(task)
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

	command := fmt.Sprintf("unzip -o %s -d /tmp/%s", params.HotFilePath, params.HotFileName)

	sshClient, err := GetSSHConn(params.ProjectId, global.OPS_CONFIG.Ops.Host, global.OPS_CONFIG.Ops.Port)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	command = strings.ReplaceAll(command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", command))
	output, err := utils.ExecuteSSHCommand(sshClient, command)
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

	command := fmt.Sprintf("bash %s /tmp/%s %s", filepath.Join(global.OPS_CONFIG.Game.GameScriptPath, "hot_game_rsync_host.sh"), params.HotFileName, params.IpList)

	sshClient, err := GetSSHConn(params.ProjectId, global.OPS_CONFIG.Ops.Host, global.OPS_CONFIG.Ops.Port)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	command = strings.ReplaceAll(command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", command))
	output, err := utils.ExecuteSSHCommand(sshClient, command)
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
	var command string
	var params NormalUpdateGameParams

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	if params.ServerType == 1 {
		command = fmt.Sprintf("bash %s %s_%d /tmp/%s/", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "hot_game_rsync_server.sh game "), params.GameType, params.GameVmid, params.HotFileName)
	} else {
		command = fmt.Sprintf("bash %s %s /tmp/%s/", filepath.Join(global.OPS_CONFIG.Game.GameScriptAutoPath, "hot_game_rsync_server.sh game_type "), params.GameType, params.HotFileName)
	}

	sshClient, err := GetSSHConn(params.ProjectId, global.OPS_CONFIG.Ops.Host, global.OPS_CONFIG.Ops.Port)
	if err != nil {
		resultList = append(resultList, "ssh连接失败")
		global.OPS_LOG.Error("ssh连接失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}
	defer sshClient.Close()

	command = strings.ReplaceAll(command, "\\", "/")

	resultList = append(resultList, fmt.Sprintf("执行命令:%s", command))
	output, err := utils.ExecuteSSHCommand(sshClient, command)
	global.OPS_LOG.Info("执行命令", zap.String("command", command), zap.String("output", output))
	resultList = append(resultList, output)
	if err != nil {
		global.OPS_LOG.Error("执行命令失败", zap.Error(err))
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, output)
	WriteTaskResult(t, resultList)
	return nil
}
