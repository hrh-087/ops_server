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

type SshAuthApi struct {
}

func (s SshAuthApi) GetSshAuthList(c *gin.Context) {
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

	list, total, err := sshAuthService.GetSshAuthList(c, info)
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

func (s SshAuthApi) CreateSshAuth(c *gin.Context) {
	var sshAuth system.SysSshAuth
	err := c.ShouldBindJSON(&sshAuth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(sshAuth, utils.SshAuthVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = sshAuthService.CreateSshAuth(c, sshAuth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

func (s SshAuthApi) UpdateSshAuth(c *gin.Context) {
	var sshAuth system.SysSshAuth
	err := c.ShouldBindJSON(&sshAuth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(sshAuth, utils.SshAuthVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = sshAuthService.UpdateSshAuth(c, sshAuth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("修改成功", c)
}
func (s SshAuthApi) DeleteSshAuth(c *gin.Context) {
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

	err = sshAuthService.DeleteSshAuth(c, id.ID)
	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func (s SshAuthApi) GetSshAuthById(c *gin.Context) {
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

	result, err := sshAuthService.GetSshAuthById(c, idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}
