package task

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"ops-server/global"
	"ops-server/model/system"
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
	resultList = append(resultList, "更新游戏镜像")
	WriteTaskResult(t, resultList)

	return nil
}

// 关闭游戏服
func HandleStopGame(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	resultList = append(resultList, "关闭游戏服")
	WriteTaskResult(t, resultList)
	return nil
}

// 更新游戏服配置
func HandleUpdateGameJsonData(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	resultList = append(resultList, "更新游戏服配置")
	WriteTaskResult(t, resultList)
	return nil
}

// 开启游戏服
func HandleStartGame(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	resultList = append(resultList, "开启游戏服")
	WriteTaskResult(t, resultList)
	return nil
}

// 检查版本号
func HandleCheckGameVersion(ctx context.Context, t *asynq.Task) error {
	var resultList []string
	resultList = append(resultList, "检查版本号")
	WriteTaskResult(t, resultList)
	return nil
}

// 热更
// 解压安装包
func HandleHotGameUnzipFile(ctx context.Context, t *asynq.Task) error {

	var params NormalUpdateGameParams
	var resultList []string

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, params.Command)
	WriteTaskResult(t, resultList)
	return nil
}

// 同步热更文件到相应服务器
func HandleHotGameRsyncHost(ctx context.Context, t *asynq.Task) error {

	var params NormalUpdateGameParams
	var resultList []string

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, params.Command)
	WriteTaskResult(t, resultList)
	return nil
}

// 同步热更文件到相应游戏服
func HandleHotGameRsyncServer(ctx context.Context, t *asynq.Task) error {
	var params NormalUpdateGameParams
	var resultList []string

	err := json.Unmarshal(t.Payload(), &params)
	if err != nil {
		resultList = append(resultList, "参数解析失败")
		WriteTaskResult(t, resultList)
		return err
	}

	resultList = append(resultList, params.Command)
	WriteTaskResult(t, resultList)
	return nil
}
