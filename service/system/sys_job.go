package system

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"time"
)

type JobService struct {
}

var JobServiceApp = new(JobService)

func (j *JobService) CreateJob(ctx *gin.Context, job system.Job) (err error) {
	claims, _ := utils.GetClaims(ctx)
	job.Creator = claims.Username
	job.CreateAt = time.Now()
	return global.OPS_DB.WithContext(ctx).Create(&job).Error
}

func (j *JobService) GetJobById(ctx context.Context, id uuid.UUID) (result system.Job, err error) {

	if err := global.OPS_DB.WithContext(ctx).Preload("Tasks").First(&result, "job_id = ?", id).Error; err == gorm.ErrRecordNotFound {
		return result, errors.New("记录不存在")
	}
	//// 获取子任务运行状态
	//if result.Status == 1 {
	//	taskNum := len(result.Tasks)
	//	completedNum := 0
	//
	//	if taskNum != 0 {
	//		for index, _ := range result.Tasks {
	//			taskInfo, err := global.AsynqInspect.GetTaskInfo("default", result.Tasks[index].AsynqId)
	//			if err != nil {
	//				//global.OPS_DB.WithContext(ctx).Delete(&result.Tasks[index])
	//				continue
	//			}
	//			result.Tasks[index].Status = taskInfo.State.String()
	//			global.OPS_DB.WithContext(ctx).Save(&result.Tasks[index])
	//
	//			if taskInfo.State.String() == "completed" || taskInfo.State.String() == "archived" {
	//				completedNum++
	//			}
	//		}
	//	}
	//
	//	// 如果子任务全部完成, 则修改作业状态为成功
	//	if taskNum == completedNum {
	//		result.Status = 2
	//		result.ExecTime = float64(time.Since(result.CreateAt).Milliseconds()) / 1000
	//		global.OPS_DB.WithContext(ctx).Save(&result)
	//	}
	//}
	return
}

func (j *JobService) GetJobList(ctx context.Context, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.Job{})

	var resultList []system.Job

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "create_at desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}
