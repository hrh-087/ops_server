package system

import (
	"errors"
	"go.uber.org/zap"
	"ops-server/global"
)

type JobTaskService struct {
}

func (j *JobTaskService) getJobTaskList() {}

func (j *JobTaskService) GetJobTaskResult(asynqId string) (result string, err error) {
	taskInfo, err := global.AsynqInspect.GetTaskInfo("default", asynqId)
	if err != nil {
		global.OPS_LOG.Error("任务结果不存在或已回收", zap.String("asynqId", asynqId), zap.Error(err))
		return "", errors.New("任务结果不存在或已回收")
	}
	return string(taskInfo.Result), err
}
