package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	sysReq "ops-server/model/system/request"
	"ops-server/utils"
	"strconv"
)

type GameUpdateApi struct {
}

func (*GameUpdateApi) CreateGameUpdate(c *gin.Context) {
	var updateParams sysReq.GameUpdateParams

	if err := c.ShouldBindJSON(&updateParams); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(updateParams.GameUpdate, utils.GameUpdateVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	updateId, err := gameUpdateService.CreateGameUpdate(c, updateParams.GameUpdate, updateParams.HotUpdateParams)
	if err != nil {
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	//response.OkWithMessage("创建成功", c)
	response.OkWithDetailed(gin.H{"id": updateId}, "创建成功", c)
}

func (*GameUpdateApi) UpdateGameUpdate(c *gin.Context) {
	var updateParams sysReq.GameUpdateParams

	if err := c.ShouldBindJSON(&updateParams); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(updateParams.GameUpdate, utils.GameUpdateVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	err := gameUpdateService.UpdateGameUpdate(c, updateParams.GameUpdate, updateParams.HotUpdateParams)
	if err != nil {
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (*GameUpdateApi) DeleteGameUpdate(c *gin.Context) {
	var ids request.GetById
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(ids, utils.IdVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	err := gameUpdateService.DeleteGameUpdate(c, ids.ID)
	if err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (*GameUpdateApi) GetGameUpdateList(c *gin.Context) {
	var info request.PageInfo
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(info, utils.PageInfoVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	result, total, err := gameUpdateService.GetGameUpdateList(c, info)
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

func (*GameUpdateApi) GetGameUpdateById(c *gin.Context) {
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

	result, err := gameUpdateService.GetGameUpdateById(c, idInfo.ID)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (*GameUpdateApi) ExecUpdateTask(c *gin.Context) {
	var idInfo request.GetById

	if err := c.ShouldBindJSON(&idInfo); err != nil {
		response.FailWithMessage("解析参数失败", c)
		return
	}

	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage("参数验证失败:"+err.Error(), c)
		return
	}

	jobId, err := gameUpdateService.ExecUpdateTask(c, idInfo.ID)
	if err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"jobId": jobId.String()}, "执行成功", c)
}

func (*GameUpdateApi) GetSvnUpdateConfigInfo(c *gin.Context) {

	result, err := gameUpdateService.GetSvnUpdateConfigInfo(c)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}
