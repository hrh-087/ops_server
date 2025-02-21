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

type GameServerApi struct{}

func (g *GameServerApi) CreateGameServer(c *gin.Context) {
	var gameServer system.SysGameServer

	err := c.ShouldBindJSON(&gameServer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(gameServer, utils.GameServerVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameServerService.CreateGameServer(c, gameServer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (g *GameServerApi) UpdateGameServer(c *gin.Context) {
	var gameServer system.SysGameServer

	err := c.ShouldBindJSON(&gameServer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(gameServer, utils.GameServerVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameServerService.UpdateGameServer(c, gameServer)
	if err != nil {
		global.OPS_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (g *GameServerApi) DeleteGameServer(c *gin.Context) {
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

	err = gameServerService.DeleteGameServer(c, id.ID)

	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)

}

func (g *GameServerApi) GetGameServerList(c *gin.Context) {
	var pageInfo systemReq.SearchGameParams

	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := gameServerService.GetGameServerList(c, pageInfo.PageInfo, pageInfo.SysGameServer)
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

func (g *GameServerApi) GetGameServerAll(c *gin.Context) {

	result, err := gameServerService.GetGameServerAll(c)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)

}

func (g *GameServerApi) GetGameServerById(c *gin.Context) {
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

	result, err := gameServerService.GetGameServerById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

func (g *GameServerApi) InstallGameServer(c *gin.Context) {
	var gameServerIds request.IdsReq

	err := c.ShouldBindJSON(&gameServerIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jobId, err := gameServerService.InstallGameServer(c, gameServerIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{
		"jobId": jobId.String(),
	}, "添加至任务至队列成功", c)
}

func (g *GameServerApi) RsyncGameConfig(c *gin.Context) {
	var updateInfo systemReq.UpdateGameConfigParams

	err := c.ShouldBindJSON(&updateInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(updateInfo, utils.UpdateGameConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	jobId, err := gameServerService.RsyncGameConfig(c, updateInfo.UpdateType, updateInfo.GameIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"jobId": jobId.String(),
	}, "添加至任务至队列成功", c)
}

func (g *GameServerApi) UpdateGameConfig(c *gin.Context) {
	var updateInfo systemReq.UpdateGameConfigParams

	err := c.ShouldBindJSON(&updateInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(updateInfo, utils.UpdateGameConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = gameServerService.UpdateGameConfig(c, updateInfo.UpdateType, updateInfo.GameIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

func (g GameServerApi) ExecGameTask(c *gin.Context) {
	var gameTask systemReq.GameTaskParams

	err := c.ShouldBindJSON(&gameTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(gameTask, utils.GameTaskVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	jobId, err := gameServerService.ExecGameTask(c, gameTask.TaskType, gameTask.GameServerIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"jobId": jobId.String(),
	}, "添加至任务至队列成功", c)

}
