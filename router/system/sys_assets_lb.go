package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AssetsLbRouter struct {
}

func (*AssetsLbRouter) InitAssetsLbRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsLbApi

	{
		router.POST("lb/rsyncCloud/", routerApi.RsyncAssetsCloudLb)
		router.POST("lb/rsyncLbListener/", routerApi.RsyncLbListener)
		router.POST("lb/writeRedis/", routerApi.WriteLBDataIntoRedis)
		router.DELETE("lb/", routerApi.DeleteAssetsLb)
		router.GET("lb/:id/", routerApi.GetAssetsLbById)
	}

	{
		routerWithoutRecord.GET("lb/", routerApi.GetAssetsLbList)
	}
}
