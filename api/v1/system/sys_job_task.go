package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/response"
	"ops-server/model/system"
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

	result, err := jobTaskService.GetJobTaskResult(jobTask.AsynqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"taskResult": result}, "获取任务结果成功", c)
}
