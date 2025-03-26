package workers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/system"
	"strings"
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

		if strings.Split(t.Type(), ":")[0] == "cron" {
			return h.ProcessTask(ctx, t)
		}

		err = json.Unmarshal(t.Payload(), &params)
		if err != nil {
			global.OPS_LOG.Error("解析taskId失败", zap.Error(err))
			return err
		}

		// 防止任务开始时,task表没有插入数据
		for i := 1; i <= 3; i++ {
			err = global.OPS_DB.First(&task, "task_id = ?", params.TaskId).Error
			if err != nil {
				if i >= 3 {
					global.OPS_LOG.Error("获取taskId失败", zap.Error(err))
					return err
				}
				time.Sleep(time.Second)
			} else {
				break
			}
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

type CronMiddlewareParams struct {
	CronTaskId uuid.UUID `json:"cronTaskId,omitempty"`
	TaskId     uuid.UUID `json:"taskId,omitempty"`
}

func CronMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		var params CronMiddlewareParams
		var err error
		var task system.JobTask
		var cronJob system.CronTask

		defer func() {
			if cronJob.Type == 1 {
				// 一次性任务只执行一次后关闭任务
				cronJob.Status = 2
				if err = global.OPS_DB.Save(&cronJob).Error; err != nil {
					global.OPS_LOG.Error("更新cronTaskId失败", zap.Error(err))
				}
			}

			if err != nil {
				_, _ = t.ResultWriter().Write([]byte(err.Error()))
				return
			}
		}()

		if err = json.Unmarshal(t.Payload(), &params); err != nil {
			global.OPS_LOG.Error("解析参数失败", zap.Error(err))
			return err
		}

		err = global.OPS_DB.First(&cronJob, "cron_task_id = ?", params.CronTaskId).Error
		if err != nil {
			global.OPS_LOG.Error("获取cronTaskId失败", zap.Error(err))
			return err
		}

		switch cronJob.Type {
		case 1:
			// 一次性任务
			err = global.OPS_DB.First(&task, "task_id = ?", params.TaskId).Error
			if err != nil {
				global.OPS_LOG.Error("获取taskId失败", zap.Error(err))
				return err
			}

		case 2:

			queue, ok := asynq.GetQueueName(ctx)
			if !ok {
				global.OPS_LOG.Error("获取队列名称失败", zap.Error(err))
				return err
			}

			// 周期性任务
			task.TaskId = uuid.Must(uuid.NewV4())
			task.JobId = cronJob.CronTaskId
			task.AsynqId = t.ResultWriter().TaskID()
			task.Status = asynq.TaskStatePending.String()
			task.HostName = global.OPS_CONFIG.Ops.Name
			task.HostIp = global.OPS_CONFIG.Ops.Host
			task.Queue = queue
			task.CreateAt = time.Now()

			if err = global.OPS_DB.Create(&task).Error; err != nil {
				global.OPS_LOG.Error("创建任务失败", zap.String("CronTaskId", cronJob.CronTaskId.String()), zap.String("taskId", task.TaskId.String()), zap.Error(err))
				return err
			}
		default:
			return errors.New("任务类型错误")
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

		// 完善task信息
		task.ExecTime = float64(execTime.Milliseconds()) / 1000

		err = global.OPS_DB.WithContext(ctx).Save(&task).Error
		if err != nil {
			return err
		}

		return nil
	})
}
