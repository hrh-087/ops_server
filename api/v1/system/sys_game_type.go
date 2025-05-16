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

type GameTypeApi struct{}

func (g *GameTypeApi) CreateGameType(c *gin.Context) {
	var gameType system.SysGameType

	err := c.ShouldBindJSON(&gameType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(gameType, utils.GameTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameTypeService.CreateGameType(c, gameType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (g *GameTypeApi) UpdateGameType(c *gin.Context) {
	var gameType system.SysGameType

	err := c.ShouldBindJSON(&gameType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(gameType, utils.GameTypeVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameTypeService.UpdateGameType(c, gameType)
	if err != nil {
		global.OPS_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (g *GameTypeApi) DeleteGameType(c *gin.Context) {
	var id request.GetById

	err := c.ShouldBindJSON(&id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(id, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameTypeService.DeleteGameType(c, id.ID)

	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)

}

func (g *GameTypeApi) GetGameTypeList(c *gin.Context) {
	var pageInfo systemReq.SearchGameTypeParams

	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := gameTypeService.GetGameTypeList(c, pageInfo.PageInfo, pageInfo.SysGameType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (g *GameTypeApi) GetGameTypeAll(c *gin.Context) {

	result, err := gameTypeService.GetGameTypeAll(c)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)

}

func (g *GameTypeApi) GetGameTypeById(c *gin.Context) {
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

	result, err := gameTypeService.GetGameTypeById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (g GameTypeApi) CopyGameType(c *gin.Context) {
	var copyGameType systemReq.CopyGameTypeParams

	if err := c.ShouldBindJSON(&copyGameType); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(copyGameType, utils.CopyGameTypeVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := gameTypeService.CopyGameType(c, copyGameType.ProjectId, copyGameType.GameTypeIds); err != nil {
		global.OPS_LOG.Error("复制失败!", zap.Error(err))
		response.FailWithMessage("复制失败"+err.Error(), c)
		return
	}

	response.OkWithMessage("复制成功", c)
	return
}
