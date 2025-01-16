package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/common/response"
	"ops-server/model/system"
	systemReq "ops-server/model/system/request"
	"ops-server/utils"
	"strconv"
)

type AssetsRedisApi struct {
}

func (a *AssetsRedisApi) CreateRedis(c *gin.Context) {
	var redis system.SysAssetsRedis

	err := c.ShouldBindJSON(&redis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(redis, utils.AssetsRedisVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = assetsRedisService.CreateRedis(c, redis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)

}

func (a *AssetsRedisApi) UpdateRedis(c *gin.Context) {
	var redis system.SysAssetsRedis

	err := c.ShouldBindJSON(&redis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(redis, utils.AssetsRedisVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = assetsRedisService.UpdateRedis(c, redis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (a *AssetsRedisApi) DeleteRedis(c *gin.Context) {
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

	err = assetsRedisService.DeleteRedis(c, idInfo.ID)

	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *AssetsRedisApi) GetAssetsRedisById(c *gin.Context) {
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

	result, err := assetsRedisService.GetAssetsRedisById(c, idInfo.ID)

	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

func (a *AssetsRedisApi) GetAssetsRedisList(c *gin.Context) {
	var pageInfo systemReq.SearchAssetsServerParams

	err := c.ShouldBindQuery(&pageInfo)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := assetsRedisService.GetAssetsRedisList(c, pageInfo.PageInfo, pageInfo.NameAndPlatformSearch)
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

func (a *AssetsRedisApi) GetAssetsRedisAll(c *gin.Context) {
	var searchInfo request.NameAndPlatformSearch

	err := c.ShouldBindQuery(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(searchInfo, utils.SearchServerVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := assetsRedisService.GetAssetsRedisAll(c, searchInfo.PlatformId)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}
