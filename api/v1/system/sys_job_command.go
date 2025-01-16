package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/model/system"
	"ops-server/utils"
	"strconv"
)

type JobCommandApi struct {
}

func (j *JobCommandApi) CreateCommand(c *gin.Context) {
	var command system.JobCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(command, utils.CommandVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := jobCommandService.CreateCommand(c, command); err != nil {
		global.OPS_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败!", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (j *JobCommandApi) UpdateCommand(c *gin.Context) {
	var command system.JobCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(command, utils.CommandVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := jobCommandService.UpdateCommand(c, command); err != nil {
		global.OPS_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败!", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (j *JobCommandApi) DeleteCommand(c *gin.Context) {
	var id request.GetById

	err := c.ShouldBindJSON(&id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(id, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := jobCommandService.DeleteCommand(c, id.ID); err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (j *JobCommandApi) GetCommandById(c *gin.Context) {
	var idInfo request.GetById
	var err error

	idStr := c.Param("id")
	idInfo.ID, err = strconv.Atoi(idStr)
	if err != nil {
		global.OPS_LOG.Error("解析id失败!", zap.Error(err))
		response.FailWithMessage("解析id失败", c)
		return
	}

	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := jobCommandService.GetCommandById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (j *JobCommandApi) GetCommandList(c *gin.Context) {
	var info request.PageInfo

	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(info, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := jobCommandService.GetCommandList(c, info)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

func (j *JobCommandApi) GetCommandAll(c *gin.Context) {
	result, err := jobCommandService.GetCommandAll(c)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (j *JobCommandApi) ExecBatchCommand(c *gin.Context) {
	var BatchInfo request.SysJobCommand

	err := c.ShouldBindJSON(&BatchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(BatchInfo, utils.BatchCommandVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	jobId, err := jobCommandService.ExecBatchCommand(c, BatchInfo)
	if err != nil {
		global.OPS_LOG.Error("创建批量执行命令任务失败!", zap.Error(err))
		response.FailWithMessage("创建批量执行命令任务失败", c)
		return
	}

	response.OkWithDetailed(gin.H{"jobId": jobId}, "执行批量执行命令任务成功", c)
}
