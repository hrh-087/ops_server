package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
)

type CronTaskService struct {
}

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
	OrderStr := "create_at desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (c CronTaskService) GetCronTaskById(ctx *gin.Context, id int) (result system.CronTask, err error) {
	if err := global.OPS_DB.WithContext(ctx).Preload("Tasks").First(&result, "job_id = ?", id).Error; err == gorm.ErrRecordNotFound {
		return result, errors.New("记录不存在")
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

func (c CronTaskService) DeleteCronTask(ctx *gin.Context, id int) (err error) {
	if err := global.OPS_DB.WithContext(ctx).Delete(&system.CronTask{}, "cron_task_id = ?", id).Error; err != nil {
		return errors.New("删除失败")
	}
	return
}

// ExecCronTask
// 开启/关闭定时任务
func (c CronTaskService) ExecCronTask(ctx *gin.Context, ids request.IdsReq) (err error) {
	return
}
