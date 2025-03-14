package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type BaseRouter struct {
}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) gin.IRouter {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi

	{
		baseRouter.POST("login/", baseApi.Login)
		baseRouter.POST("logout/", baseApi.Logout)
		baseRouter.POST("uploadFile/", baseApi.UploadFile)
		baseRouter.GET("generateExcel/", baseApi.GenerateExcel)
	}
	return baseRouter
}
