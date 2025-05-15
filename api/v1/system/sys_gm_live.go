package system

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/utils"
)

func (g GmApi) GetLiveServerList(c *gin.Context) {
	var live request.GmServerId
	if err := c.ShouldBindJSON(&live); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(live, utils.GmServerIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := gmService.GetActivityServerList(c, live.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (g GmApi) ClearLiveServer(c *gin.Context) {
	var live request.GmServerId
	if err := c.ShouldBindJSON(&live); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(live, utils.GmServerIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := gmService.ClearActivityServer(c, live.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("清除成功", c)
}

func (g GmApi) SetLiveActivityServer(c *gin.Context) {
	var setLiveData request.LiveServerActivity
	if err := c.ShouldBindJSON(&setLiveData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(setLiveData, utils.SetLiveActivityServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := gmService.SetLiveActivityServer(c, setLiveData)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("设置成功", c)

}

func (g GmApi) KickLivePlayer(c *gin.Context) {
	var kickLiveData request.LiveKickPlayer
	if err := c.ShouldBindJSON(&kickLiveData); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(kickLiveData, utils.GmServerIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := gmService.KickLivePlayer(c, kickLiveData)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
}
