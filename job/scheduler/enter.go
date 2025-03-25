package scheduler

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"log"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/job/workers"
	"ops-server/model/system"
	"time"
)

func NewAsynqScheduler() *asynq.Scheduler {

	loc, err := time.LoadLocation(global.OPS_CONFIG.Asynq.Tz)
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}

	return asynq.NewScheduler(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", global.OPS_CONFIG.Asynq.Addr, global.OPS_CONFIG.Asynq.Port),
		Password: global.OPS_CONFIG.Asynq.Pass,
		DB:       global.OPS_CONFIG.Asynq.Db,
	}, &asynq.SchedulerOpts{
		Location: loc,
	})
}

func InitScheduler() {
	global.AsynqScheduler = NewAsynqScheduler()

	if global.AsynqScheduler == nil {
		panic("获取Scheduler失败")
	}

	var cronTaskList []system.CronTask
	// 获取已开启的周期性任务初始化
	err := global.OPS_DB.Where("type = 2 and status = 1").Find(&cronTaskList).Error
	if err != nil {

		panic("获取定时任务列表失败")
	}

	err = global.OPS_DB.Transaction(func(tx *gorm.DB) error {
		for index := range cronTaskList {
			entryId, err := task.NewCronTask(cronTaskList[index].TaskTypeName, workers.CronMiddlewareParams{
				CronTaskId: cronTaskList[index].CronTaskId,
			}, cronTaskList[index].CronRule, asynq.Queue("cron"))

			if err != nil {
				panic("添加定时任务失败")
			}

			cronTaskList[index].TaskId = entryId

			if err := tx.Save(&cronTaskList[index]).Error; err != nil {
				panic("更新定时任务失败")
			}
		}

		return nil
	})

	if err != nil {
		panic("初始化定时任务失败")
	}

	if err := global.AsynqScheduler.Run(); err != nil {
		panic("启动定时任务失败")
	}
}
