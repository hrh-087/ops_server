package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/utils"
)

func (g GmApi) GetPlayerId(c *gin.Context) {
	var importPlayerParams request.ImportPlayerDataParams

	if err := c.ShouldBindJSON(&importPlayerParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := gmService.GetPlayerIdByAccountName(c, importPlayerParams.Account, importPlayerParams.PlatformCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (g GmApi) ImportPlayerData(c *gin.Context) {
	var importPlayerParams request.ImportPlayerDataParams
	if err := c.ShouldBindJSON(&importPlayerParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(importPlayerParams, utils.GmImportPlayerDataVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := gmService.ImportPlayerData(c, importPlayerParams.Account, importPlayerParams.ImportId, importPlayerParams.OutputPlayerId, importPlayerParams.PlatformCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("导入成功", c)
}
