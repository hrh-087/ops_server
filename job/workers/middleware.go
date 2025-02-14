package workers

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"time"
)

type ExecTimeMiddlewareParams struct {
	TaskId uuid.UUID
}

func GetExecTimeMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		var params ExecTimeMiddlewareParams
		var err error
		var task system.JobTask
		var job system.Job

		err = json.Unmarshal(t.Payload(), &params)
		if err != nil {
			global.OPS_LOG.Error("解析taskId失败", zap.Error(err))
			return err
		}

		startTime := time.Now()
		err = h.ProcessTask(ctx, t)
		if err != nil {
			global.OPS_LOG.Error("任务执行失败", zap.Error(err), zap.String("taskId", task.TaskId.String()))
			task.Status = asynq.TaskStateArchived.String()
		} else {
			task.Status = asynq.TaskStateCompleted.String()
		}
		execTime := time.Since(startTime)

		err = global.OPS_DB.First(&task, "task_id = ?", params.TaskId).Error
		if err != nil {
			global.OPS_LOG.Error("获取taskId失败", zap.Error(err))
			return err
		}

		// 完善task信息
		task.ExecTime = float64(execTime.Milliseconds()) / 1000

		err = global.OPS_DB.WithContext(ctx).Save(&task).Error
		if err != nil {
			return err
		}

		// 检查作业任务状态
		err = global.OPS_DB.WithContext(ctx).Preload("Tasks").First(&job, "job_id = ?", task.JobId).Error
		if err != nil {
			global.OPS_LOG.Error("获取jobId失败", zap.Error(err), zap.String("job_id", task.JobId.String()))
		}

		taskNum := len(job.Tasks)
		completedNum := 0
		failNum := 0

		if taskNum != 0 {
			for _, task := range job.Tasks {
				if task.Status == "completed" {
					completedNum++
				} else if task.Status == "archived" {
					failNum++
				}
			}
		}

		if taskNum == (completedNum + failNum) {
			if failNum == 0 {
				job.Status = 2
			} else {
				job.Status = 3
			}
			job.ExecTime = float64(time.Since(job.CreateAt).Milliseconds()) / 1000
			err = global.OPS_DB.WithContext(ctx).Save(&job).Error
		}

		return nil
	})
}
