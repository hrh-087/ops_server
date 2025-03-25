package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/job/workers"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"regexp"
	"sync"
	"time"
)

type CronTaskService struct {
}

const CronRegex = "^(\\*|([0-5]?\\d)|(\\*/[1-5]?\\d))\\s+(\\*|([01]?\\d|2[0-3])|(\\*/\\d+))\\s+(\\*|([1-9]|[12]\\d|3[01])|(\\*/\\d+))\\s+(\\*|([1-9]|1[0-2])|(\\*/\\d+))\\s+(\\*|([0-7])|(\\*/\\d+))$"

func (c CronTaskService) GetCronTaskList(ctx *gin.Context, info request.PageInfo, cronTask system.CronTask) (result interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.CronTask{})

	var resultList []system.CronTask

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if cronTask.Status != 0 {
		db = db.Where("status = ?", cronTask.Status)
	}

	if cronTask.Type != 0 {
		db = db.Where("type = ?", cronTask.Type)
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "created_at desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (c CronTaskService) GetCronTaskById(ctx *gin.Context, id uuid.UUID) (result system.CronTask, err error) {
	if err := global.OPS_DB.WithContext(ctx).Preload("Tasks").First(&result, "cron_task_id = ?", id).Error; err == gorm.ErrRecordNotFound {
		return result, errors.New("记录不存在")
	} else if err != nil {
		return result, err
	}
	return
}

func (c CronTaskService) CreateCronTask(ctx *gin.Context, cronTask system.CronTask) (err error) {
	claims, _ := utils.GetClaims(ctx)
	cronTask.Creator = claims.Username
	cronTask.CronTaskId = uuid.Must(uuid.NewV4())
	return global.OPS_DB.WithContext(ctx).Create(&cronTask).Error
}

func (c CronTaskService) UpdateCronTask(ctx *gin.Context, cronTask system.CronTask) (err error) {
	var old system.CronTask

	if old.Status == 1 {
		return errors.New("任务正在运行, 请先停止任务")
	}

	updateField := []string{
		"Describe",
		"Name",
	}

	if errors.Is(global.OPS_DB.WithContext(ctx).Where("cron_task_id = ?", cronTask.CronTaskId).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	return global.OPS_DB.WithContext(ctx).Model(&old).Select(updateField).Updates(cronTask).Error
}

func (c CronTaskService) DeleteCronTask(ctx *gin.Context, id uuid.UUID) (err error) {
	var old system.CronTask
	if errors.Is(global.OPS_DB.WithContext(ctx).Where("cron_task_id = ?", id).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	if old.Status == 1 {
		return errors.New("任务正在运行, 请先停止任务")
	}

	if err := global.OPS_DB.WithContext(ctx).Delete(&system.CronTask{}, "cron_task_id = ?", id).Error; err != nil {
		return errors.New("删除失败")
	}
	return
}

// ExecCronTask
// 开启/关闭定时任务
// 一次性任务跟周期性任务 的异步任务不通用
func (c CronTaskService) ExecCronTask(ctx *gin.Context, cronTask system.CronTask) (err error) {
	var old system.CronTask
	var mu sync.Mutex
	var t system.JobTask

	mu.Lock()
	defer mu.Unlock()

	err = global.OPS_DB.WithContext(ctx).Where("cron_task_id = ?", cronTask.CronTaskId).First(&old).Error
	if err != nil {
		return err
	}

	switch old.Status {
	case 1:
		if old.TaskId == "" {
			return errors.New("TaskId为空, 请先开启任务")
		}
		if old.Type == 2 {
			// todo 重启服务的话需要重载定时任务
			err = global.AsynqScheduler.Unregister(old.TaskId)
			if err != nil {
				global.OPS_LOG.Error("关闭周期性定时任务失败", zap.Error(err))
				return err
			}
		} else if old.Type == 1 {
			//if old.ExecTime.Before(time.Now()) {
			//	return errors.New("任务已执行")
			//}
			err = global.AsynqInspect.DeleteTask("cron", old.TaskId)
			if err != nil {
				global.OPS_LOG.Error("关闭一次性定时任务失败", zap.Error(err))
				return err
			}
		}
		old.Status = 2
	case 2:

		taskId := uuid.Must(uuid.NewV4())
		if old.Type == 2 {
			if cronTask.CronRule == "" {
				return errors.New("cron规则为空")
			}
			match, _ := regexp.MatchString(CronRegex, cronTask.CronRule)
			if !match {
				return errors.New("cron规则格式错误")
			}

			//asynqTask := task.NewTask(old.TaskTypeName, []byte{})
			entryId, err := task.NewCronTask(old.TaskTypeName, workers.CronMiddlewareParams{
				CronTaskId: old.CronTaskId,
			}, cronTask.CronRule, asynq.Queue("cron"))
			if err != nil {
				global.OPS_LOG.Error("开启周期性定时任务失败", zap.Error(err))
				return err
			}
			old.CronRule = cronTask.CronRule
			old.TaskId = entryId

		} else if old.Type == 1 {
			if cronTask.ExecTime.IsZero() {
				return errors.New("执行时间不能为空")
			}
			taskInfo, err := task.NewOnceTask(old.TaskTypeName, workers.CronMiddlewareParams{
				CronTaskId: old.CronTaskId,
				TaskId:     taskId,
			}, asynq.ProcessAt(cronTask.ExecTime), asynq.Queue("cron"))
			if err != nil {
				global.OPS_LOG.Error("开启一次性定时任务失败", zap.Error(err))
				return err
			}

			t.TaskId = taskId
			t.JobId = old.CronTaskId
			t.AsynqId = taskInfo.ID
			t.Status = taskInfo.State.String()
			t.HostName = global.OPS_CONFIG.Ops.Name
			t.HostIp = global.OPS_CONFIG.Ops.Host
			t.CreateAt = time.Now()

			if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
				global.OPS_LOG.Error("创建任务失败", zap.String("CronTaskId", old.CronTaskId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
				return err
			}

			old.ExecTime = cronTask.ExecTime
			old.TaskId = taskInfo.ID

		}
		old.Status = 1
	}
	if err := global.OPS_DB.WithContext(ctx).Model(&old).Updates(old).Error; err != nil {
		global.OPS_LOG.Error("更新定时任务失败", zap.Error(err))
		return err
	}

	return
}
