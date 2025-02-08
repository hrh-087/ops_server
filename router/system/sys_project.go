package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type ProjectRouter struct {
}

func (p *ProjectRouter) InitProjectRouter(Router *gin.RouterGroup) {
	projectRouter := Router.Group("project").Use(middleware.OperationRecord())
	projectRouterWithoutRecord := Router.Group("project")

	projectRouterApi := v1.ApiGroupApp.SystemApiGroup.ProjectApi

	{
		projectRouter.POST("createProject/", projectRouterApi.CreateProject)
		projectRouter.POST("updateProject/", projectRouterApi.UpdateProject)
		projectRouter.DELETE("deleteProject/", projectRouterApi.DeleteProject)
		projectRouter.POST("getProjectById/", projectRouterApi.GetProjectById)
		projectRouter.POST("setAuthorityProject/", projectRouterApi.SetAuthorityProject)
	}
	{
		projectRouterWithoutRecord.POST("getProjectList/", projectRouterApi.GetProjectList)
		projectRouterWithoutRecord.GET("getBaseProjectTree/", projectRouterApi.GetBaseProjectTree)
		projectRouterWithoutRecord.POST("getAuthorityProject/", projectRouterApi.GetAuthorityProject)
	}

}
