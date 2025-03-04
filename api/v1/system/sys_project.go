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
)

type ProjectApi struct{}

// CreateProject
// @Tags      SysProjectApi
// @Summary   创建项目
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data body      system.SysProject                true  "创建项目"
// @Success   200   {object}  response.Response{msg=string}  "创建项目"
// @Router    /project/createProject [post]
func (s *ProjectApi) CreateProject(c *gin.Context) {
	var project system.SysProject

	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(project, utils.ProjectVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = projectService.CreateProject(project)
	if err != nil {
		global.OPS_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateProject
// @Tags      SysProjectApi
// @Summary   修改项目
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data body      system.SysProject                true  "创建项目"
// @Success   200   {object}  response.Response{msg=string}  "创建项目"
// @Router    /project/updateProject [put]
func (s *ProjectApi) UpdateProject(c *gin.Context) {
	var project system.SysProject

	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(project, utils.ProjectVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = projectService.UpdateProject(project)
	if err != nil {
		global.OPS_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败"+err.Error(), c)
		return
	}

	response.OkWithMessage("修改成功", c)
}

func (s *ProjectApi) GetProjectById(c *gin.Context) {
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
	result, err := projectService.GetProjectById(idInfo.ID)
	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

func (s *ProjectApi) GetProjectList(c *gin.Context) {
	var pageInfo systemReq.SearchProjectParams

	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := projectService.GetProjectList(pageInfo.SysProject, pageInfo.PageInfo)
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

func (s *ProjectApi) DeleteProject(c *gin.Context) {
	var project system.SysProject
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(project.OpsModel, utils.IdVerify)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = projectService.DeleteProject(project)
	if err != nil {
		global.OPS_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (s *ProjectApi) GetBaseProjectTree(c *gin.Context) {
	projects, err := projectService.GetAllProject()

	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(projects, "获取成功", c)
}

func (s *ProjectApi) GetAuthorityProject(c *gin.Context) {
	var param request.GetAuthorityId
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	projects, err := projectService.GetAuthorityProject(&param)

	if err != nil {
		global.OPS_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(projects, "获取成功", c)
}

func (s *ProjectApi) SetAuthorityProject(c *gin.Context) {
	var authorityProject systemReq.AddAuthorityProject
	err := c.ShouldBindJSON(&authorityProject)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(authorityProject, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := projectService.SetAuthorityProject(authorityProject.AuthorityId, authorityProject.ProjectIds); err != nil {
		global.OPS_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

func (s *ProjectApi) InitProject(c *gin.Context) {
	var project system.SysProject
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := utils.Verify(project, utils.ProjectVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	_, err = projectService.InitProject(c, project)
	if err != nil {
		global.OPS_LOG.Error("初始化失败!", zap.Error(err))
		response.FailWithMessage("初始化失败"+err.Error(), c)
		return
	}

	response.OkWithMessage("已添加至初始化队列", c)
}
