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

func (g GmApi) GetRankList(c *gin.Context) {
	var rankParams request.GmRankOpenParams
	if err := c.ShouldBindJSON(&rankParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(rankParams, utils.GmRankOpenVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := gmService.GetRankList(c, rankParams.ServerId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "获取成功", c)
}

func (g GmApi) GetRankRewardList(c *gin.Context) {
	var rankParams request.GmRankRewardParams
	if err := c.ShouldBindJSON(&rankParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(rankParams, utils.GmRankRewardVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := gmService.GetRankRewardList(c, rankParams.ServerId, rankParams.Id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)

}

func (g GmApi) SetRankConfig(c *gin.Context) {
	var rankParams request.GmRankConfigParams
	if err := c.ShouldBindJSON(&rankParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(rankParams, utils.GmRankConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := gmService.SetRankConfig(c, rankParams.ServerId, rankParams.RankConfig); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("设置成功", c)
}

func (g GmApi) UploadGameConfig(c *gin.Context) {
	var uploadParams request.UploadFileParams
	if err := c.ShouldBindJSON(&uploadParams); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(uploadParams, utils.FileNameVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := gmService.UploadGameConfig(c, uploadParams.FileName); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("上传成功", c)
}
