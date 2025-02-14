package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"time"
)

type SysTaskService struct {
}

func (s *SysTaskService) CreateTask(ctx *gin.Context, task system.SysTask) (err error) {
	return global.OPS_DB.WithContext(ctx).Create(&task).Error
}

func (s *SysTaskService) UpdateTask(ctx *gin.Context, task system.SysTask) (err error) {
	return global.OPS_DB.WithContext(ctx).Save(&task).Error
}

func (s *SysTaskService) DeleteTask(ctx *gin.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&system.SysTask{}).Error
}

func (s *SysTaskService) GetTaskList(ctx *gin.Context, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysTask{})

	var resultList []system.SysTask
	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (s *SysTaskService) GetTaskById(ctx *gin.Context, id int) (result system.SysTask, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return
}

func (s *SysTaskService) ExecTask(ctx *gin.Context, id int) (jobId uuid.UUID, err error) {
	var job system.Job
	var taskManage system.SysTask
	var taskList []system.JobTask
	var host system.SysAssetsServer
	var t system.JobTask
	if err = global.OPS_DB.WithContext(ctx).Where("id = ?", id).First(&taskManage).Error; err != nil {
		return
	}

	if err = global.OPS_DB.WithContext(ctx).Where("server_type = 3").First(&host).Error; err != nil {
		return
	}

	jobId = uuid.Must(uuid.NewV4())
	taskId := uuid.Must(uuid.NewV4())

	taskInfo, err := task.NewCommonTask(taskManage.TaskType, task.CommonTaskParams{
		ProjectId: taskManage.ProjectId,
		TaskId:    taskId,
		HostId:    host.ID,
	})

	if err != nil {
		global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}

	t.JobId = jobId
	t.AsynqId = taskInfo.ID
	t.TaskId = taskId
	t.Status = taskInfo.State.String()
	t.HostName = host.ServerName
	t.HostIp = host.PubIp
	t.CreateAt = time.Now()

	if err = global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
		global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
		return
	}

	taskList = append(taskList, t)

	job.JobId = jobId
	job.Name = taskManage.Name
	job.Status = 1
	job.Type = taskManage.TaskType
	job.Tasks = taskList

	// 创建作业任务
	err = JobServiceApp.CreateJob(ctx, job)
	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}

	return
}
