package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/response"
	"ops-server/model/system"
	"ops-server/model/system/request"
	"ops-server/utils"
)

type JobTaskApi struct {
}

func (j *JobTaskApi) GetJobTaskResult(c *gin.Context) {
	var jobTask system.JobTask

	err := c.ShouldBindJSON(&jobTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := jobTaskService.GetJobTaskResult(jobTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"taskResult": result}, "获取任务结果成功", c)
}

func (j *JobTaskApi) GetJobTaskList(c *gin.Context) {

	var jobTaskInfo request.JobTaskParams
	err := c.ShouldBindQuery(&jobTaskInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(jobTaskInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := jobTaskService.GetJobTaskList(c, jobTaskInfo.PageInfo, jobTaskInfo.JobId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     jobTaskInfo.Page,
		PageSize: jobTaskInfo.PageSize,
	}, "获取成功", c)
}
