package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type JobTaskService struct {
}

var JobTaskServiceApp = new(JobTaskService)

func (j *JobTaskService) GetJobTaskList(ctx *gin.Context, info request.PageInfo, jobId string) (result interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	if jobId == "" {
		return nil, 0, errors.New("jobId不能为空")
	}

	db := global.OPS_DB.WithContext(ctx).Model(&system.JobTask{}).Where("job_id = ?", jobId)

	var resultList []system.JobTask

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Error
	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "create_at desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (j *JobTaskService) GetJobTaskResult(jobTask system.JobTask) (result string, err error) {
	taskInfo, err := global.AsynqInspect.GetTaskInfo(jobTask.Queue, jobTask.AsynqId)
	if err != nil {
		global.OPS_LOG.Error("任务结果不存在或已回收", zap.String("asynqId", jobTask.AsynqId), zap.String("queue", jobTask.Queue), zap.Error(err))
		return "", errors.New("任务结果不存在或已回收")
	}
	return string(taskInfo.Result), err
}

// CancelJobTask 关闭一次性定时任务
func (j *JobTaskService) CancelJobTask(queue, asynqId string) (err error) {
	err = global.AsynqInspect.DeleteTask(queue, asynqId)
	if err != nil {
		global.OPS_LOG.Error("关闭一次性定时任务失败", zap.Error(err))
		return err
	}

	return global.OPS_DB.Model(&system.JobTask{}).Where("asynq_id = ?", asynqId).Update("status", "cancel").Error
}
