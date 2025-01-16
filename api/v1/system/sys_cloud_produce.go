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

type CloudProduceApi struct{}

func (s *CloudProduceApi) CreateCloudProduce(c *gin.Context) {
	var cloud system.SysCloudProduce

	err := c.ShouldBindJSON(&cloud)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(cloud, utils.CloudVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = cloudProduceService.CreateCloudProduce(c, cloud)

	if err != nil {
		global.OPS_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (s CloudProduceApi) GetCloudProduceList(c *gin.Context) {
	var pageInfo systemReq.SearchCloudProduceParams

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

	list, total, err := cloudProduceService.GetCloudProduceList(c, pageInfo.SysCloudProduce, pageInfo.PageInfo)

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

func (s CloudProduceApi) GetCloudProduceById(c *gin.Context) {
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

	result, err := cloudProduceService.GetCloudProduceById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (s CloudProduceApi) UpdateCloudProduce(c *gin.Context) {
	var cloudProduce system.SysCloudProduce

	err := c.ShouldBindJSON(&cloudProduce)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(cloudProduce, utils.CloudVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = cloudProduceService.UpdateCloudProduce(c, cloudProduce)
	if err != nil {
		global.OPS_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (s CloudProduceApi) DeleteCloudProduce(c *gin.Context) {
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

	err = cloudProduceService.DeleteCloudProduce(c, idInfo.ID)

	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (s *CloudProduceApi) GetCloudProduceAll(c *gin.Context) {

	result, err := cloudProduceService.GetCloudProduceAll(c)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}