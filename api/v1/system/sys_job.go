package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/model/system"
	"ops-server/utils"
)

type JobApi struct {
}

func (j *JobApi) GetJobList(c *gin.Context) {
	var page request.PageInfo

	err := c.ShouldBindQuery(&page)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(page, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := jobService.GetJobList(c, page)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, "获取成功", c)

}

func (j *JobApi) GetJobById(c *gin.Context) {
	var job system.Job

	err := c.ShouldBindJSON(&job)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(job, utils.JobVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	job, err = jobService.GetJobById(c, job.JobId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(job, "获取成功", c)
}
