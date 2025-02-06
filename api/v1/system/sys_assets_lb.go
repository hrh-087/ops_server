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

type AssetsLbApi struct {
}

func (a *AssetsLbApi) RsyncAssetsCloudLb(c *gin.Context) {
	var lb system.SysAssetsLb
	err := c.ShouldBindJSON(&lb)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(lb, utils.AssetsLbVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = assetsLbService.RsyncAssetsCloudLb(c, lb)
	if err != nil {
		response.FailWithMessage("同步华为云负载均衡列表失败", c)
		return
	}

	response.OkWithMessage("同步华为云负载均衡列表成功", c)
}

func (a *AssetsLbApi) RsyncLbListener(c *gin.Context) {
	err := assetsLbService.RsyncLbListener(c)
	if err != nil {
		response.FailWithMessage("同步华为云负载均衡监听器失败", c)
		return
	}

	response.OkWithMessage("同步华为云负载均衡监听器成功", c)
}

func (a *AssetsLbApi) WriteLBDataIntoRedis(c *gin.Context) {
	err := assetsLbService.WriteLBDataIntoRedis(c)
	if err != nil {
		response.FailWithMessage("同步redis失败", c)
		return
	}
	response.OkWithMessage("同步redis成功", c)
}

func (a *AssetsLbApi) GetAssetsLbList(c *gin.Context) {
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

	list, total, err := assetsLbService.GetAssetsLbList(c, info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

func (a *AssetsLbApi) DeleteAssetsLb(c *gin.Context) {
	var ids request.GetById
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(ids, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = assetsLbService.DeleteAssetsLb(c, ids.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (a *AssetsLbApi) GetAssetsLbById(c *gin.Context) {
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

	result, err := assetsLbService.GetAssetsLbById(c, idInfo.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
