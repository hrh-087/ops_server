package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type BaseRouter struct {
}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base").Use(middleware.ProjectAuth())
	baseRouterWithOutProject := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi

	{
		baseRouter.GET("generateExcel/", baseApi.GenerateExcel)
	}

	{
		baseRouterWithOutProject.POST("login/", baseApi.Login)
		baseRouterWithOutProject.POST("logout/", baseApi.Logout)
		baseRouterWithOutProject.POST("uploadFile/", baseApi.UploadFile)
	}
	//return baseRouter
}
