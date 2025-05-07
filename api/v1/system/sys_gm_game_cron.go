package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/utils"
)

func (g GmApi) GetGameCronList(c *gin.Context) {
	var params request.GameServerActivityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := gmService.GetGameCronList(c, params.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "获取成功", c)
}

func (g GmApi) GetActivityExtra(c *gin.Context) {
	var params request.GameServerActivityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := gmService.GetActivityExtra(c, params.ServerId, params.Key)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "获取成功", c)
}

func (g GmApi) SetGameCron(c *gin.Context) {
	var params request.GameServerActivityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(params, utils.GameServerActivityParamsVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := gmService.SetGameCron(c, params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("设置成功", c)
}
