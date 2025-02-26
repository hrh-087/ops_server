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

type GamePlatformApi struct{}

func (g *GamePlatformApi) CreatePlatform(c *gin.Context) {
	var platform system.SysGamePlatform

	err := c.ShouldBindJSON(&platform)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(platform, utils.PlatformVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = platformService.CreatePlatform(c, platform)
	if err != nil {
		global.OPS_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (g *GamePlatformApi) UpdatePlatform(c *gin.Context) {
	var platform system.SysGamePlatform

	err := c.ShouldBindJSON(&platform)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(platform, utils.PlatformVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = platformService.UpdatePlatform(c, platform)
	if err != nil {
		global.OPS_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (g *GamePlatformApi) DeletePlatform(c *gin.Context) {
	var idInfo request.GetById

	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = platformService.DeletePlatform(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)

}

func (g *GamePlatformApi) GetPlatformList(c *gin.Context) {
	var pageInfo request.PageInfo

	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := platformService.GetPlatformList(c, system.SysGamePlatform{}, pageInfo)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)

}

func (g *GamePlatformApi) GetPlatformById(c *gin.Context) {
	var idInfo request.GetById
	var err error

	idStr := c.Param("id")
	idInfo.ID, err = strconv.Atoi(idStr)

	if err != nil {
		global.OPS_LOG.Error("解析id失败!", zap.Error(err))
		response.FailWithMessage("解析id失败", c)
		return
	}

	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := platformService.GetPlatformById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

func (g *GamePlatformApi) GetPlatformAll(c *gin.Context) {

	result, err := platformService.GetPlatformAll(c)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (g *GamePlatformApi) KickGameServer(c *gin.Context) {
	var platform system.SysGamePlatform
	err := c.ShouldBindJSON(&platform)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	serverId, err := strconv.ParseInt(platform.PlatformCode, 10, 64)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = platformService.KickGameServer(c, int(serverId))
	if err != nil {
		global.OPS_LOG.Error("踢人失败!", zap.Error(err))
		response.FailWithMessage("踢人失败", c)
		return
	}

	response.OkWithMessage("踢人成功", c)
}
