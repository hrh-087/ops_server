package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AssetsServerRouter struct {
}

func (a *AssetsServerRouter) InitAssetsServerRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsServerApi

	{
		router.GET("server/:id/", routerApi.GetAssetsServerById)
		router.POST("server/", routerApi.CreateServer)
		router.PUT("server/", routerApi.UpdateServer)
		router.DELETE("server/", routerApi.DeleteServer)
		router.POST("server/generatePrometheusHostConfig/", routerApi.GeneratePrometheusHostConfig)
		router.POST("server/pullInstanceCloudInfo/", routerApi.PullInstanceCloudInfo)
	}
	{
		routerWithoutRecord.GET("server/all/", routerApi.GetAssetsServerAll)
		routerWithoutRecord.GET("server/", routerApi.GetAssetsServerList)
	}
}
