package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ops-server/global"
	request2 "ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/model/system"
	"ops-server/model/system/request"

	"ops-server/utils"
	"strconv"
)

type CronTaskApi struct {
}

func (CronTaskApi) GetCronTaskList(c *gin.Context) {
	var info request.CronTaskParams
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(info, utils.PageInfoVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	result, total, err := cronTaskService.GetCronTaskList(c, info.PageInfo, info.CronTask)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     result,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)

}

func (CronTaskApi) GetCronTaskById(c *gin.Context) {
	var idInfo request2.GetById
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

	result, err := cronTaskService.GetCronTaskById(c, idInfo.ID)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (CronTaskApi) CreateCronTask(c *gin.Context) {
	var cronTask system.CronTask

	if err := c.ShouldBindJSON(&cronTask); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(cronTask, utils.CronTaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := cronTaskService.CreateCronTask(c, cronTask); err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (CronTaskApi) UpdateCronTask(c *gin.Context) {

}

func (CronTaskApi) DeleteCronTask(c *gin.Context) {
	var ids request2.GetById
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(ids, utils.IdVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	if err := cronTaskService.DeleteCronTask(c, ids.ID); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
