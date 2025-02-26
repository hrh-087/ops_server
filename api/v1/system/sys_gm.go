package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/utils"
)

type GmApi struct {
}

func (g GmApi) SetSwitch(c *gin.Context) {
	var gmParams request.GmSwitchParams
	if err := c.ShouldBindJSON(&gmParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(gmParams, utils.GmSwitchVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	msg, err := gmService.SetSwitch(c, gmParams.ServerId, gmParams.TypeKey, gmParams.State)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(msg, c)
}

func (g GmApi) GetSwitchList(c *gin.Context) {
	var gmParams request.GmSwitchParams
	if err := c.ShouldBindJSON(&gmParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := gmService.GetSwitchList(c, gmParams.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "获取成功", c)
}
