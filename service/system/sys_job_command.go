package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"time"
)

type JobCommandService struct {
}

func (j *JobCommandService) CreateCommand(ctx *gin.Context, command system.JobCommand) (err error) {
	err = global.OPS_DB.WithContext(ctx).Create(&command).Error
	return
}

func (j *JobCommandService) UpdateCommand(ctx *gin.Context, command system.JobCommand) (err error) {
	var old system.JobCommand
	updateField := []string{
		"Name",
		"Command",
		"CommandType",
		"Description",
		"UseBatch",
	}

	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", command.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	err = global.OPS_DB.WithContext(ctx).Model(&old).Select(updateField).Updates(command).Error
	return
}

func (j *JobCommandService) DeleteCommand(ctx *gin.Context, id int) (err error) {
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.JobCommand{}).Error
	return
}

func (j *JobCommandService) GetCommandById(ctx *gin.Context, id int) (result system.JobCommand, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return
}

func (j *JobCommandService) GetCommandList(ctx *gin.Context, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.JobCommand{})

	var resultList []system.JobCommand

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"

	err = db.Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (j *JobCommandService) GetCommandAll(ctx *gin.Context) (result []system.JobCommand, err error) {
	err = global.OPS_DB.WithContext(ctx).Where("use_batch = ?", true).Find(&result).Error
	return
}

func (j *JobCommandService) ExecBatchCommand(ctx *gin.Context, batchInfo request.SysJobCommand) (jobId uuid.UUID, err error) {
	var job system.Job
	var taskList []system.JobTask
	var hostServerList []system.SysAssetsServer
	var commandInfo system.JobCommand

	switch batchInfo.BatchType {
	case 1:
		err = global.OPS_DB.WithContext(ctx).Where("id in ?", batchInfo.ServerList).Find(&hostServerList).Error
		if err != nil {
			return uuid.UUID{}, errors.New("获取服务器列表失败")
		}
	case 2:
		err = global.OPS_DB.WithContext(ctx).Where("platform_id in ?", batchInfo.ServerList).Find(&hostServerList).Error
		if err != nil {
			return uuid.UUID{}, errors.New("获取服务器列表失败")
		}
	default:
		return uuid.UUID{}, errors.New("未知类型")
	}

	err = global.OPS_DB.WithContext(ctx).Where("id = ?", batchInfo.CommandId).First(&commandInfo).Error
	if err != nil {
		return uuid.UUID{}, errors.New("获取命令失败")
	}

	// 创建任务
	jobId = uuid.Must(uuid.NewV4())
	for index, _ := range hostServerList {
		var t system.JobTask

		taskId := uuid.Must(uuid.NewV4())
		taskInfo, err := task.NewBatchCommand(task.BatchCommand{
			Command: commandInfo.Command,
			Host:    hostServerList[index],
			TaskId:  taskId,
		})
		if err != nil {
			global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
			continue
		}

		t.JobId = jobId
		t.AsynqId = taskInfo.ID
		t.TaskId = taskId
		t.Status = taskInfo.State.String()
		t.HostName = hostServerList[index].ServerName
		t.HostIp = hostServerList[index].PubIp
		t.CreateAt = time.Now()

		if err := global.OPS_DB.WithContext(ctx).Create(&t).Error; err != nil {
			global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
			continue
		}
		taskList = append(taskList, t)
	}
	job.JobId = jobId
	job.Name = "批量执行命令"
	job.Status = 1
	job.Type = "batchCommand"
	job.Tasks = taskList

	err = JobServiceApp.CreateJob(ctx, job)
	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}
	return jobId, nil
}
